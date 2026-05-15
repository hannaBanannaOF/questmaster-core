package middleware

import (
	"questmaster-core/internal/shared/context"

	"github.com/gin-gonic/gin"
)

func QueryParamsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		params := make(map[string]string)

		queryValues := c.Request.URL.Query()

		for key, values := range queryValues {
			if len(values) > 0 {
				params[key] = values[0]
			}
		}

		appCtx := context.AppContext{Context: c}
		appCtx.SetFilters(params)

		// Segue para o próximo handler
		c.Next()
	}
}
