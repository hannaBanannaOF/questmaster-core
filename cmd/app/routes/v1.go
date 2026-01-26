package routes

import (
	campaignTransport "questmaster-core/internal/transport/campaign/http"
	characterTransport "questmaster-core/internal/transport/character/http"

	"github.com/gin-gonic/gin"
)

type V1RoutesDeps struct {
	CampaignHandler  *campaignTransport.CampaignsHandler
	CharacterHandler *characterTransport.CharactersHandler
	AuthMiddleware   gin.HandlerFunc
	PermMiddleware   gin.HandlerFunc
}

func RegisterV1Routes(router *gin.Engine, deps V1RoutesDeps) {
	v1 := router.Group(
		"/core/api/v1",
		deps.AuthMiddleware,
		deps.PermMiddleware,
	)

	registerCampaignRoutes(v1, deps.CampaignHandler)
	registerCharacterRoutes(v1, deps.CharacterHandler)
}

// v1 := router.Group("/core/api/v1", rpg.AuthMiddleware(rsa), rpg.UpdatePermissionsMiddleware(mongoClient, pgPool))
// {
//me := v1.Group("/me")
//{
//me.GET("/campaings", handlers.MyCampaingsHandler(q := queries.ListUserCampaignsQuery()))
// me.GET("/calendar", meHdlrV1.GetMyCalendar)
// me.GET("/upcoming", meHdlrV1.GetMyUpcoming)
// me.GET("/sheets", meHdlrV1.GetMyCharacterSheets)
//}
// campaign := v1.Group("/campaign")
// {
// campaign.GET("", campaignHandler.GetMyCampaigns)
// 	campaing.POST("", campaingHdlrV1.CreateCampaing)
// 	campaing.GET("/resolve/:slug", perm.CheckViewCampaingPermission(pgPool), campaingHdlrV1.ResolveSlug)
// 	campaing.GET("/:campaingId", perm.CheckViewCampaingPermission(pgPool), campaingHdlrV1.GetCampaingDetails)
// 	campaing.PUT("/:campaingId/toggle-in-play", perm.CheckUpdateCampaingPermission(pgPool, func(ctx *gin.Context) (int, error) {
// 		return strconv.Atoi(ctx.Param("campaingId"))
// 	}), campaingHdlrV1.ToggleCampaingInPlay)
// 	campaing.POST("/:campaingId/schedule", perm.CheckUpdateCampaingPermission(pgPool, func(ctx *gin.Context) (int, error) {
// 		return strconv.Atoi(ctx.Param("campaingId"))
// 	}), campaingHdlrV1.ScheduleSession)
// }
// characterSheet := v1.Group("/character-sheet")
// {
// 	characterSheet.POST("", characterSheetHdlrV1.CreateCharacterSheet)
// 	characterSheet.GET("/resolve/:slug", perm.CheckViewCharacterSheetPermission(pgPool), characterSheetHdlrV1.ResolveSlug)
// 	characterSheet.GET("/:characterSheetId", perm.CheckViewCharacterSheetPermission(pgPool), characterSheetHdlrV1.GetCharacterSheetDetails)
// }
// invite := v1.Group("/invite")
// {
// 	invite.GET("/:inviteHash", inviteHldrV1.GetInviteDetails)
// 	invite.POST("/:inviteHash/accept", inviteHldrV1.AcceptInvite)
// 	invite.POST("/create", perm.CheckUpdateCampaingPermission(pgPool, func(ctx *gin.Context) (int, error) {
// 		var body vo.CreateInvite
// 		if err := ctx.ShouldBindBodyWith(&body, binding.JSON); err != nil {
// 			return 0, err
// 		}
// 		return body.CampaingId, nil
// 	}), inviteHldrV1.GetCampaingInvite)
// }
// }
