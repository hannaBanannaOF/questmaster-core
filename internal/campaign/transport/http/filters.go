package campaign

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"

	campaignDomain "questmaster-core/internal/campaign/domain"
	"questmaster-core/internal/shared/context"
	"questmaster-core/internal/shared/httperrors"
	"questmaster-core/internal/shared/pagination"
)

// CampaignListFiltersMiddleware parses the role, status, q and without_players filters of the campaign list
func CampaignListFiltersMiddleware() gin.HandlerFunc {
	return campaignFiltersMiddleware(true)
}

// CampaignRoleFilterMiddleware parses only the role filter, for routes where the list filters don't apply
func CampaignRoleFilterMiddleware() gin.HandlerFunc {
	return campaignFiltersMiddleware(false)
}

func campaignFiltersMiddleware(list bool) gin.HandlerFunc {
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

		if raw := c.Query("status"); list && raw != "" {
			status, err := campaignDomain.NewCampaignStatus(raw)
			if err != nil {
				_ = c.Error(err)
				c.Abort()
				return
			}
			filters.Status = &status
		}

		if list {
			search, err := pagination.NewSearch(c.Query("q"))
			if err != nil {
				_ = c.Error(fmt.Errorf("%w: %w", httperrors.ErrInvalidParam, err))
				c.Abort()
				return
			}
			filters.Search = search

			if raw := c.Query("without_players"); raw != "" {
				withoutPlayers, err := strconv.ParseBool(raw)
				if err != nil {
					_ = c.Error(fmt.Errorf("%w: without_players", httperrors.ErrInvalidParam))
					c.Abort()
					return
				}
				filters.WithoutPlayers = withoutPlayers
			}
		}

		appCtx := context.AppContext{Context: c}
		appCtx.SetCampaignListFilters(filters)

		c.Next()
	}
}
