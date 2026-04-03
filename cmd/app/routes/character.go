package routes

import (
	characterTransport "questmaster-core/internal/character/transport/http"
	appContext "questmaster-core/internal/shared/context"
	middleware "questmaster-core/internal/shared/middleware"

	"github.com/gin-gonic/gin"
)

func registerCharacterRoutes(
	v1 *gin.RouterGroup,
	handler *characterTransport.CharactersHandler,
) {
	character := v1.Group("/character")
	{
		character.GET("", appContext.Adapt(handler.GetCurrentUserCharacters))
		character.POST("", appContext.Adapt(handler.CreateCharacter))
		character.GET("/resolve/:slug", middleware.Slug(), appContext.Adapt(handler.ResolveSlug))
		character.GET("/:characterID", middleware.CharacterID(), appContext.Adapt(handler.GetDetails))
		character.PATCH("/:characterID/hp", middleware.CharacterID(), appContext.Adapt(handler.UpdateHP))
		character.DELETE("/:characterID", middleware.CharacterID(), appContext.Adapt(handler.DeleteCharacter))
	}
}
