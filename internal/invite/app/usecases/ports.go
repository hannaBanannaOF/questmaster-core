package invite

import (
	campaignDomain "questmaster-core/internal/campaign/domain"
	characterDomain "questmaster-core/internal/character/domain"
	rpgDomain "questmaster-core/internal/rpg/domain"
)

type InviteCampaignFinder interface {
	FindByID(campaignID campaignDomain.CampaignID) (campaignDomain.Campaign, error)
}

type InviteAvailableCharacterFinder interface {
	GetBySystemAndCampaignIDNull(userID rpgDomain.UserID, system rpgDomain.System) ([]characterDomain.Character, error)
}

type InviteCharacterCampaignLinker interface {
	LinkToCampaign(campaignID campaignDomain.CampaignID, characterID characterDomain.CharacterID) (characterDomain.Character, error)
}
