package character

import campaignDomain "questmaster-core/internal/campaign/domain"

type CharacterCampaingFinder interface {
	FindByID(campaignID campaignDomain.CampaignID) (campaignDomain.Campaign, error)
}
