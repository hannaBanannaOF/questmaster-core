package invite

import (
	campaignDomain "questmaster-core/internal/campaign/domain"
	inviteDomain "questmaster-core/internal/invite/domain"
)

type InviteRepository interface {
	Create(campaignID campaignDomain.CampaignID) (*inviteDomain.Invite, error)
	FindByCampaignID(campaignID campaignDomain.CampaignID) (*inviteDomain.Invite, error)
	FindByHash(hash inviteDomain.InviteHash) (*inviteDomain.Invite, error)
}
