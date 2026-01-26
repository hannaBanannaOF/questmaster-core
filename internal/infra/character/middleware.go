package character

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CheckViewCharacterSheetPermission(pg *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("characterSheetId"))
		if err != nil {
			c.Next()
			return
		}
		uid := c.GetString("UserID")
		ctx := context.Background()
		rows := pg.QueryRow(ctx, "select c.id from character_sheet c left join campaign s on s.id = c.campaign_id WHERE (s.dm_id = $1 OR c.player_id = $1) AND c.id = $2", uid, id)
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
