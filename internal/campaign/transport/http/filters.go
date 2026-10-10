package campaign

import (
	"github.com/gin-gonic/gin"

	campaignDomain "questmaster-core/internal/campaign/domain"
	"questmaster-core/internal/shared/context"
)

// CampaignListFiltersMiddleware parses the role and status filters of the campaign list
func CampaignListFiltersMiddleware() gin.HandlerFunc {
	return campaignFiltersMiddleware(true)
}

// CampaignRoleFilterMiddleware parses only the role filter, for routes where status doesn't apply
func CampaignRoleFilterMiddleware() gin.HandlerFunc {
	return campaignFiltersMiddleware(false)
}

func campaignFiltersMiddleware(withStatus bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var filters campaignDomain.CampaignListFilters

		if raw := c.Query("role"); raw != "" {
			role, err := campaignDomain.NewCampaignRole(raw)
			if err != nil {
				_ = c.Error(err)
				c.Abort()
				return
			}
			filters.Role = &role
		}

		if raw := c.Query("status"); withStatus && raw != "" {
			status, err := campaignDomain.NewCampaignStatus(raw)
			if err != nil {
				_ = c.Error(err)
				c.Abort()
				return
			}
			filters.Status = &status
		}

		appCtx := context.AppContext{Context: c}
		appCtx.SetCampaignListFilters(filters)

		c.Next()
	}
}
