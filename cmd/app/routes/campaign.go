package routes

import (
	campaignTransport "questmaster-core/internal/campaign/transport/http"
	rpgTransport "questmaster-core/internal/rpg/transport/http"
	appContext "questmaster-core/internal/shared/context"
	"questmaster-core/internal/shared/middleware"

	"github.com/gin-gonic/gin"
)

func registerCampaignRoutes(
	v1 *gin.RouterGroup,
	handler *campaignTransport.CampaignHandler,
) {
	campaign := v1.Group("/campaign")
	{
		campaign.GET("", middleware.PageMiddleware(), campaignTransport.CampaignListFiltersMiddleware(), appContext.Adapt(handler.GetCurrentUserCampaigns))
		campaign.POST("", appContext.Adapt(handler.CreateCampaign))
		campaign.GET("/resolve/:slug", rpgTransport.SlugMiddleware(), appContext.Adapt(handler.ResolveSlug))
		campaign.DELETE("/:campaignID", campaignTransport.CampaignIDMiddleware(), appContext.Adapt(handler.DeleteCampaign))
		campaign.GET("/:campaignID", campaignTransport.CampaignIDMiddleware(), appContext.Adapt(handler.GetCampaignDetails))
		campaign.PATCH("/:campaignID/status", campaignTransport.CampaignIDMiddleware(), appContext.Adapt(handler.UpdateStatus))
	}
}
