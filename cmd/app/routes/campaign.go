package routes

import (
	campaignTransport "questmaster-core/internal/campaign/transport/http"
	appContext "questmaster-core/internal/shared/context"
	middleware "questmaster-core/internal/shared/middleware"

	"github.com/gin-gonic/gin"
)

func registerCampaignRoutes(
	v1 *gin.RouterGroup,
	handler *campaignTransport.CampaignHandler,
) {
	campaign := v1.Group("/campaign")
	{
		campaign.GET("", appContext.Adapt(handler.GetMyCampaigns))
		campaign.POST("", appContext.Adapt(handler.CreateCampaign))
		campaign.GET("/resolve/:slug", middleware.Slug(), appContext.Adapt(handler.ResolveSlug))
		campaign.DELETE("/:campaignID", middleware.CampaignID(), appContext.Adapt(handler.DeleteCampaign))
		campaign.GET("/:campaignID/details", middleware.CampaignID(), appContext.Adapt(handler.GetCampaignDetails))
		campaign.PATCH("/:campaignID/status", middleware.CampaignID(), appContext.Adapt(handler.UpdateStatus))
		campaign.POST("/:campaignID/invite", middleware.CampaignID(), appContext.Adapt(handler.GetOrCreateCampaignInvite))
	}
}
