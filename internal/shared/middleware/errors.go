package middleware

import (
	"log"

	"github.com/gin-gonic/gin"

	"questmaster-core/internal/shared/httperrors"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

		// Executa handlers seguintes
		c.Next()

		// Se não houve erro, segue o fluxo normal
		if len(c.Errors) == 0 {
			return
		}

		// Pega o último erro registrado
		lastErr := c.Errors.Last().Err

		httpErr := httperrors.From(lastErr)

		// Log interno
		log.Printf(
			"[HTTP ERROR] %d %s %s -> %v",
			httpErr.Status,
			c.Request.Method,
			c.Request.URL.Path,
			lastErr,
		)

		// Evita escrever duas vezes
		if c.Writer.Written() {
			return
		}

		c.AbortWithStatusJSON(
			httpErr.Status,
			httpErr,
		)
	}
}
