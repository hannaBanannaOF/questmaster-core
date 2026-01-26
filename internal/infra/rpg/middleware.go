package rpg

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func AuthMiddleware(jwkSet map[string]*rsa.PublicKey) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Expecting "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		accessToken := parts[1]

		token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			// Extract the key ID (kid) from the token header
			kid, ok := token.Header["kid"].(string)
			if !ok {
				return nil, fmt.Errorf("kid not found in token header")
			}

			// Retrieve the public key from the JWK set using the key ID
			publicKey, keyExists := jwkSet[kid]
			if !keyExists {
				return nil, fmt.Errorf("public key not found for kid: %v", kid)
			}

			return publicKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		sub, err := token.Claims.GetSubject()
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("UserID", sub)
		c.Next()
	}
}

func GetJWKSet(url string) (map[string]*rsa.PublicKey, error) {
	// Make the GET request
	response, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error making GET request: %v", err)
	}
	defer response.Body.Close()

	// Decode the JSON response
	var jwkSet struct {
		Keys []struct {
			Kid string   `json:"kid"`
			N   string   `json:"n"`
			E   string   `json:"e"`
			X5C []string `json:"x5c"`
		} `json:"keys"`
	}
	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(&jwkSet); err != nil {
		return nil, fmt.Errorf("error decoding JSON: %v", err)
	}

	// Create a map to store RSA public keys
	jwkMap := make(map[string]*rsa.PublicKey)

	// Iterate through each key in the JWK set
	for _, key := range jwkSet.Keys {
		// Decode base64url-encoded modulus (N) and exponent (E)
		modulus, err := decodeBase64URL(key.N)
		if err != nil {
			return nil, fmt.Errorf("error decoding modulus: %v", err)
		}

		exponent, err := decodeBase64URL(key.E)
		if err != nil {
			return nil, fmt.Errorf("error decoding exponent: %v", err)
		}

		// Create RSA public key
		pubKey := &rsa.PublicKey{
			N: modulus,
			E: int(exponent.Int64()),
		}

		// Store the public key in the map using the key ID (Kid)
		jwkMap[key.Kid] = pubKey
	}

	return jwkMap, nil
}

// decodeBase64URL decodes a base64url-encoded string and returns a big.Int
func decodeBase64URL(input string) (*big.Int, error) {
	// Convert base64url to base64
	base64Str := strings.ReplaceAll(input, "-", "+")
	base64Str = strings.ReplaceAll(base64Str, "_", "/")

	// Pad the base64 string with "="
	switch len(base64Str) % 4 {
	case 2:
		base64Str += "=="
	case 3:
		base64Str += "="
	}

	// Decode base64 string
	data, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		return nil, err
	}

	// Convert bytes to big.Int
	result := new(big.Int).SetBytes(data)
	return result, nil
}

func UpdatePermissionsMiddleware(mongo *mongo.Client, pg *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("UserID")
		if userID != "" {
			go func(uid string) {
				ctx := context.Background()

				fichas := []int{}
				rows, err := pg.Query(ctx, "select c.id from character_sheet c left join campaign s on s.id = c.campaign_id WHERE s.dm_id = $1 OR c.player_id = $1", uid)
				if err != nil {
					log.Printf("Unable to get user permitted sheets: %s\n", err)
				}
				defer rows.Close()
				for rows.Next() {
					var s int
					err := rows.Scan(&s)
					if err != nil {
						log.Printf("Unable to get user permitted sheets: %s\n", err)
					}
					fichas = append(fichas, s)
				}

				mesas := []int{}
				rows, err = pg.Query(ctx, "SELECT DISTINCT s.id FROM campaign s LEFT JOIN character_sheet c ON c.campaign_id = s.id WHERE s.dm_id = $1 OR c.player_id = $1", uid)
				if err != nil {
					log.Printf("Unable to get user permitted campaigns: %s\n", err)
				}
				for rows.Next() {
					var s int
					err := rows.Scan(&s)
					if err != nil {
						log.Printf("Unable to get user permitted campaigns: %s\n", err)
					}
					mesas = append(mesas, s)
				}
				filter := bson.M{"player_id": uid}
				update := bson.M{
					"$set": bson.M{
						"fichas":     fichas,
						"campanhas":  mesas,
						"updated_at": time.Now(),
					},
				}

				_, err = mongo.Database(os.Getenv("MONGODB_DATABASENAME")).Collection("permissions").
					UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))

				if err != nil {
					log.Printf("Erro ao atualizar permissões do usuário %s: %v", uid, err)
				}
			}(userID)
		}

		c.Next()
	}
}
