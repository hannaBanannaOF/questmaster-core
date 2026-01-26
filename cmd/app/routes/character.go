package routes

import (
	transport "questmaster-core/internal/transport/character/http"

	"github.com/gin-gonic/gin"
)

func registerCharacterRoutes(
	v1 *gin.RouterGroup,
	handler *transport.CharactersHandler,
) {
	character := v1.Group("/character")
	{
		character.GET("", handler.GetMyCharacters)
		character.POST("", handler.CreateCharacter)
		character.GET("/resolve/:slug", handler.ResolveSlug)
		character.GET("/:characterId", handler.GetDetails)
	}
}
