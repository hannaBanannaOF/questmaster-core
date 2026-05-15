package character

import rpgDomain "questmaster-core/internal/rpg/domain"

type CharacterListFilters struct {
	GameSystem      *rpgDomain.System
	WithoutCampaign *bool
}
