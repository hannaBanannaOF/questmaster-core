package character

import (
	campaignDomain "questmaster-core/internal/campaign/domain"
	characterDomain "questmaster-core/internal/character/domain"
	rpgDomain "questmaster-core/internal/rpg/domain"
	"questmaster-core/internal/shared/pagination"
	userDomain "questmaster-core/internal/user/domain"
)

type CreateCharacterCommand struct {
	Name   characterDomain.CharacterName
	Hp     *characterDomain.HP
	System rpgDomain.System
	Player userDomain.UserID
}

type CreatePremadeCharacterCommand struct {
	CampaignID campaignDomain.CampaignID
	Name       characterDomain.CharacterName
	Hp         *characterDomain.HP
	UserID     userDomain.UserID
}

type EditCharacterCommand struct {
	ID     characterDomain.CharacterID
	Name   *characterDomain.CharacterName
	MaxHP  *int
	UserID userDomain.UserID
}

type LeaveCampaignCommand struct {
	CampaignID  campaignDomain.CampaignID
	CharacterID characterDomain.CharacterID
	UserID      userDomain.UserID
}

type DeleteCharacterCommand struct {
	ID     characterDomain.CharacterID
	UserID userDomain.UserID
}

type GetCharacterDetailsCommand struct {
	ID     characterDomain.CharacterID
	UserID userDomain.UserID
}

type GetCurrentUserCharactersCommand struct {
	UserID  userDomain.UserID
	Filters characterDomain.CharacterListFilters
	Page    pagination.Page
}

type UpdateHPCommand struct {
	ID     characterDomain.CharacterID
	NewHP  int
	UserID userDomain.UserID
}
