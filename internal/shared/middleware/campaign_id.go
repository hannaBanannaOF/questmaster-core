package middleware

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"

	campaignDomain "questmaster-core/internal/campaign/domain"
	"questmaster-core/internal/shared/context"
	"questmaster-core/internal/shared/httperrors"
)

func CampaignID() gin.HandlerFunc {
	return func(c *gin.Context) {

		raw := c.Param("campaignID")

		id, err := strconv.Atoi(raw)
		if err != nil {
			_ = c.Error(
				fmt.Errorf("%w: campaignID", httperrors.ErrInvalidParam),
			)
			c.Abort()
			return
		}

		appCtx := context.AppContext{Context: c}
		appCtx.SetCampaignID(campaignDomain.NewCampaignID(id))

		c.Next()
	}
}
