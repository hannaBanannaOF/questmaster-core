package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"questmaster-core/internal/shared/context"
	"questmaster-core/internal/shared/httperrors"
	"questmaster-core/internal/shared/pagination"
)

// PageMiddleware parses the limit and offset query params of a list route
func PageMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		page, err := pagination.NewPage(c.Query("limit"), c.Query("offset"))
		if err != nil {
			_ = c.Error(fmt.Errorf("%w: %w", httperrors.ErrInvalidParam, err))
			c.Abort()
			return
		}

		appCtx := context.AppContext{Context: c}
		appCtx.SetPage(page)

		c.Next()
	}
}
