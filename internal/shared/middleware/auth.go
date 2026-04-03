package middleware

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"questmaster-core/internal/shared/context"
	"questmaster-core/internal/shared/httperrors"
	userDomain "questmaster-core/internal/user/domain"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func AuthMiddleware(jwkSet map[string]*rsa.PublicKey) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			_ = c.Error(httperrors.ErrUnauthorized)
			c.Abort()
			return
		}

		// Expecting "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			_ = c.Error(httperrors.ErrUnauthorized)
			c.Abort()
			return
		}

		accessToken := parts[1]

		token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
			if token.Method.Alg() != jwt.SigningMethodRS256.Alg() {
				return nil, fmt.Errorf("unexpected alg")
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
			_ = c.Error(httperrors.ErrUnauthorized)
			c.Abort()
			return
		}

		sub, err := token.Claims.GetSubject()
		if err != nil {
			_ = c.Error(httperrors.ErrUnauthorized)
			c.Abort()
			return
		}

		uuid, err := uuid.Parse(sub)
		if err != nil {
			_ = c.Error(httperrors.ErrUnauthorized)
			c.Abort()
			return
		}

		claimsMap := token.Claims.(jwt.MapClaims)

		firstName := getOptionalString(claimsMap, "given_name")
		lastName := getOptionalString(claimsMap, "family_name")
		var compositeName *userDomain.Name
		if firstName != nil {
			o, err := userDomain.NewName(getString(claimsMap, "given_name"), lastName)
			if err != nil {
				_ = c.Error(httperrors.ErrUnauthorized)
				c.Abort()
				return
			}
			compositeName = &o
		}

		username, err := userDomain.NewUsername(getString(claimsMap, "preferred_username"))
		if err != nil {
			_ = c.Error(httperrors.ErrUnauthorized)
			c.Abort()
			return
		}

		user := userDomain.User{
			Id:       userDomain.NewUserID(uuid),
			Username: username,
			Name:     compositeName,
		}

		appCtx := context.AppContext{Context: c}
		appCtx.SetUser(user)

		c.Next()
	}
}

func getString(claims jwt.MapClaims, key string) string {
	if val, ok := claims[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

func getOptionalString(claims jwt.MapClaims, key string) *string {
	if val, ok := claims[key]; ok {
		if str, ok := val.(string); ok && str != "" {
			return &str
		}
	}
	return nil
}

type OIDCDiscovery struct {
	Issuer  string `json:"issuer"`
	JWKSUri string `json:"jwks_uri"`
	// outros campos disponíveis se precisar futuramente
	AuthorizationEndpoint string `json:"authorization_endpoint"`
	TokenEndpoint         string `json:"token_endpoint"`
	UserinfoEndpoint      string `json:"userinfo_endpoint"`
}

func GetJWKSet(oidcHost string) (map[string]*rsa.PublicKey, error) {
	wellKnownURL := strings.TrimRight(oidcHost, "/") + "/.well-known/openid-configuration"

	resp, err := http.Get(wellKnownURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch OIDC discovery document: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("discovery endpoint returned status %d", resp.StatusCode)
	}

	var discovery OIDCDiscovery
	if err := json.NewDecoder(resp.Body).Decode(&discovery); err != nil {
		return nil, fmt.Errorf("failed to decode discovery document: %w", err)
	}

	if discovery.JWKSUri == "" {
		return nil, fmt.Errorf("jwks_uri not found in discovery document")
	}
	// Make the GET request
	response, err := http.Get(discovery.JWKSUri)
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
