package character

import (
	"strconv"

	"github.com/gin-gonic/gin"

	characterDomain "questmaster-core/internal/character/domain"
	rpgDomain "questmaster-core/internal/rpg/domain"
	"questmaster-core/internal/shared/context"
)

// CharacterListFiltersMiddleware parses the game_system and without_campaign filters of the character list
func CharacterListFiltersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var filters characterDomain.CharacterListFilters

		if raw := c.Query("game_system"); raw != "" {
			system, err := rpgDomain.NewSystem(raw)
			if err != nil {
				_ = c.Error(err)
				c.Abort()
				return
			}
			filters.GameSystem = &system
		}

		// Kept lenient like before: values that aren't booleans don't filter
		if withoutCampaign, err := strconv.ParseBool(c.Query("without_campaign")); err == nil {
			filters.WithoutCampaign = &withoutCampaign
		}

		appCtx := context.AppContext{Context: c}
		appCtx.SetCharacterListFilters(filters)

		c.Next()
	}
}
