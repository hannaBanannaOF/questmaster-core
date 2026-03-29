package middleware

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"

	characterDomain "questmaster-core/internal/character/domain"
	"questmaster-core/internal/shared/context"
	"questmaster-core/internal/shared/httperrors"
)

func CharacterID() gin.HandlerFunc {
	return func(c *gin.Context) {

		raw := c.Param("characterID")

		id, err := strconv.Atoi(raw)
		if err != nil {
			_ = c.Error(
				fmt.Errorf("%w: characterID", httperrors.ErrInvalidParam),
			)
			c.Abort()
			return
		}

		appCtx := context.AppContext{Context: c}
		appCtx.SetCharacterID(characterDomain.NewCharacterID(id))

		c.Next()
	}
}
