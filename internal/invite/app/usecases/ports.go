package invite

import (
	campaignDomain "questmaster-core/internal/campaign/domain"
	characterDomain "questmaster-core/internal/character/domain"
	rpgDomain "questmaster-core/internal/rpg/domain"
	userDomain "questmaster-core/internal/user/domain"
)

type InviteCampaignFinder interface {
	FindByID(campaignID campaignDomain.CampaignID) (campaignDomain.Campaign, error)
}

type InviteCharacterCampaignLinker interface {
	LinkToCampaign(campaignID campaignDomain.CampaignID, characterSlug rpgDomain.Slug, userID userDomain.UserID) (characterDomain.Character, error)
}
