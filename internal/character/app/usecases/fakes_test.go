package character

import campaignDomain "questmaster-core/internal/campaign/domain"

// fakeCampaignFinder serves one campaign to the use cases that load the character campaign
type fakeCampaignFinder struct {
	campaign campaignDomain.Campaign
}

func (f fakeCampaignFinder) FindByID(campaignDomain.CampaignID) (campaignDomain.Campaign, error) {
	return f.campaign, nil
}
