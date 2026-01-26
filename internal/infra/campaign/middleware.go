package campaign

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CheckCampaignId func(*gin.Context) (int, error)

func CheckViewCampaignPermission(pg *pgxpool.Pool, check CheckCampaignId) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := check(c)
		if err != nil {
			c.Next()
			return
		}
		uid := c.GetString("UserID")
		ctx := context.Background()
		rows := pg.QueryRow(ctx, "SELECT s.id FROM campaign s LEFT JOIN character_sheet c ON c.campaign_id = s.id WHERE (s.dm_id = $1 OR c.player_id = $1) AND s.id = $2", uid, id)
		var s int
		err = rows.Scan(&s)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				c.Next()
				return
			}
			c.JSON(http.StatusForbidden, gin.H{"error": "Can't view resource!"})
			c.Abort()
		}
		c.Next()
	}
}
