package character

import (
	campaignDomain "questmaster-core/internal/campaign/domain"
	characterDomain "questmaster-core/internal/character/domain"
	rpgDomain "questmaster-core/internal/rpg/domain"
)

type CharacterRepository interface {
	GetAllByPlayerID(userID rpgDomain.UserID) ([]characterDomain.Character, error)
	GetAllByCampaignID(campaignID campaignDomain.CampaignID) ([]characterDomain.Character, error)
	FindBySlug(slug rpgDomain.Slug) (*characterDomain.Character, error)
	FindByID(characterID characterDomain.CharacterID) (*characterDomain.Character, error)
	Create(name characterDomain.CharacterName, playerID rpgDomain.UserID, system rpgDomain.System, hp *characterDomain.HP) (characterDomain.Character, error)
	UpdateHP(newHP characterDomain.HP, characterID characterDomain.CharacterID) (characterDomain.Character, error)
	DeleteByID(characterID characterDomain.CharacterID) (bool, error)
	GetAllByUserIDAndCampaignIDNullAndSystem(userID rpgDomain.UserID, system rpgDomain.System) ([]characterDomain.Character, error)
	UpdateCampaign(campaignID campaignDomain.CampaignID, characterID characterDomain.CharacterID) (*characterDomain.Character, error)
}
