package invite

import campaignDomain "questmaster-core/internal/campaign/domain"

type Invite struct {
	Id         InviteID
	CampaignId campaignDomain.CampaignID
	Hash       InviteHash
}
