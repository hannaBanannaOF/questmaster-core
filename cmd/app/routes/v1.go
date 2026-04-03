package routes

import (
	campaignTransport "questmaster-core/internal/campaign/transport/http"
	characterTransport "questmaster-core/internal/character/transport/http"
	inviteTransport "questmaster-core/internal/invite/transport/http"
	userTransport "questmaster-core/internal/user/transport/http"

	"github.com/gin-gonic/gin"
)

type V1RoutesDeps struct {
	CampaignHandler  *campaignTransport.CampaignHandler
	CharacterHandler *characterTransport.CharactersHandler
	InviteHandler    *inviteTransport.InviteHandler
	UserHandler      *userTransport.UserHandler
	AuthMiddleware   gin.HandlerFunc
	PermMiddleware   gin.HandlerFunc
}

func RegisterV1Routes(router *gin.Engine, deps V1RoutesDeps) {
	v1 := router.Group(
		"/core/api/v1",
		deps.AuthMiddleware,
		deps.PermMiddleware,
	)

	registerUserRoutes(v1, deps.UserHandler)
	registerCampaignRoutes(v1, deps.CampaignHandler)
	registerCharacterRoutes(v1, deps.CharacterHandler)
	registerInviteRoutes(v1, deps.InviteHandler)
}
