package routes

import (
	appContext "questmaster-core/internal/shared/context"
	userTransport "questmaster-core/internal/user/transport/http"

	"github.com/gin-gonic/gin"
)

func registerUserRoutes(
	v1 *gin.RouterGroup,
	handler *userTransport.UserHandler,
) {
	user := v1.Group("/user")
	{
		user.GET("", appContext.Adapt(handler.GetInfo))
	}
}
