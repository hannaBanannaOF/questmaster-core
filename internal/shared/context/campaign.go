package context

import campaignDomain "questmaster-core/internal/campaign/domain"

func (c *AppContext) SetCampaignID(id campaignDomain.CampaignID) {
	c.Set(string(campaignIDKey), id)
}

func (c *AppContext) CampaignID() campaignDomain.CampaignID {
	v, ok := c.Get(string(campaignIDKey))
	if !ok {
		panic("CampaignID not found in context")
	}

	id, ok := v.(campaignDomain.CampaignID)
	if !ok {
		panic("CampaignID has invalid type")
	}

	return id
}
