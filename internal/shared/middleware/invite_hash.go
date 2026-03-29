package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	inviteDomain "questmaster-core/internal/invite/domain"
	"questmaster-core/internal/shared/context"
	"questmaster-core/internal/shared/httperrors"
)

func InviteHash() gin.HandlerFunc {
	return func(c *gin.Context) {

		raw := c.Param("inviteHash")

		hash, err := uuid.Parse(raw)
		if err != nil {
			_ = c.Error(
				fmt.Errorf("%w: inviteHash", httperrors.ErrInvalidParam),
			)
			c.Abort()
			return
		}

		appCtx := context.AppContext{Context: c}
		appCtx.SetInviteHash(inviteDomain.NewHash(hash))

		c.Next()
	}
}
