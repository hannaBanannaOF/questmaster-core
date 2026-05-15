package routes

import (
	characterTransport "questmaster-core/internal/character/transport/http"
	rpgTransport "questmaster-core/internal/rpg/transport/http"
	appContext "questmaster-core/internal/shared/context"

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
		character.GET("/resolve/:slug", rpgTransport.SlugMiddleware(), appContext.Adapt(handler.ResolveSlug))
		character.GET("/:characterID", characterTransport.CharacterIDMiddleware(), appContext.Adapt(handler.GetDetails))
		character.PATCH("/:characterID/hp", characterTransport.CharacterIDMiddleware(), appContext.Adapt(handler.UpdateHP))
		character.DELETE("/:characterID", characterTransport.CharacterIDMiddleware(), appContext.Adapt(handler.DeleteCharacter))
	}
}
