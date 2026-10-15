package routes

import (
	campaignTransport "questmaster-core/internal/campaign/transport/http"
	characterTransport "questmaster-core/internal/character/transport/http"
	rpgTransport "questmaster-core/internal/rpg/transport/http"
	appContext "questmaster-core/internal/shared/context"
	"questmaster-core/internal/shared/middleware"

	"github.com/gin-gonic/gin"
)

func registerCharacterRoutes(
	v1 *gin.RouterGroup,
	handler *characterTransport.CharactersHandler,
) {
	character := v1.Group("/character")
	{
		character.GET("", middleware.PageMiddleware(), characterTransport.CharacterListFiltersMiddleware(), appContext.Adapt(handler.GetCurrentUserCharacters))
		character.POST("", appContext.Adapt(handler.CreateCharacter))
		character.GET("/resolve/:slug", rpgTransport.SlugMiddleware(), appContext.Adapt(handler.ResolveSlug))
		character.GET("/:characterID", characterTransport.CharacterIDMiddleware(), appContext.Adapt(handler.GetDetails))
		character.PATCH("/:characterID", characterTransport.CharacterIDMiddleware(), appContext.Adapt(handler.Edit))
		character.DELETE("/:characterID", characterTransport.CharacterIDMiddleware(), appContext.Adapt(handler.DeleteCharacter))
	}

	// Premade characters live in the character module but are created inside a campaign
	v1.POST("/campaign/:campaignID/character", campaignTransport.CampaignIDMiddleware(), appContext.Adapt(handler.CreatePremadeCharacter))
	v1.DELETE("/campaign/:campaignID/character/:characterID", campaignTransport.CampaignIDMiddleware(), characterTransport.CharacterIDMiddleware(), appContext.Adapt(handler.RemoveFromCampaign))
}
