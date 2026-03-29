package middleware

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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
