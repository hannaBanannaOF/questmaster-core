package routes

import (
	transport "questmaster-core/internal/transport/campaign/http"

	"github.com/gin-gonic/gin"
)

func registerCampaignRoutes(
	v1 *gin.RouterGroup,
	handler *transport.CampaignsHandler,
) {
	campaign := v1.Group("/campaign")
	{
		campaign.GET("", handler.GetMyCampaigns)
		campaign.POST("", handler.CreateCampaign)
		campaign.GET("/resolve/:slug", handler.ResolveSlug)
		campaign.GET("/:campaignId", handler.GetCampaignDetails)
		campaign.DELETE("/:campaignId", handler.DeleteCampaign)
		campaign.PATCH("/:campaignId/status", handler.UpdateStatus)
	}
}
