package routes

import (
	inviteTransport "questmaster-core/internal/invite/transport/http"
	appContext "questmaster-core/internal/shared/context"

	"github.com/gin-gonic/gin"
)

func registerInviteRoutes(
	v1 *gin.RouterGroup,
	handler *inviteTransport.InviteHandler,
) {
	invite := v1.Group("/invite")
	{
		invite.POST("", appContext.Adapt(handler.CreateInvite))
		invite.GET("/:inviteHash", inviteTransport.InviteHashMiddleware(), appContext.Adapt(handler.GetInviteDetails))
		invite.POST("/:inviteHash/accept", inviteTransport.InviteHashMiddleware(), appContext.Adapt(handler.AcceptInvite))
	}
}
