package rpg

import (
	"github.com/gin-gonic/gin"

	rpg "questmaster-core/internal/rpg/domain"
	"questmaster-core/internal/shared/context"
)

func SlugMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		raw := c.Param("slug")

		appCtx := context.AppContext{Context: c}

		slug, err := rpg.NewSlug(raw)
		if err != nil {
			_ = c.Error(err)
			c.Abort()
			return
		}
		appCtx.SetSlug(slug)

		c.Next()
	}
}
