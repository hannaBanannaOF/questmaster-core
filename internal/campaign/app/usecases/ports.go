package campaign

import (
	campaignDomain "questmaster-core/internal/campaign/domain"
	characterDomain "questmaster-core/internal/character/domain"
	inviteDomain "questmaster-core/internal/invite/domain"
)

type CampaignCharacterFinder interface {
	GetByCampaignID(campaignID campaignDomain.CampaignID) ([]characterDomain.Character, error)
}

type CampaignInviteFinder interface {
	GetByCampaignID(campaignID campaignDomain.CampaignID) (inviteDomain.Invite, error)
}

type CampaignInviteCreator interface {
	Create(campaignID campaignDomain.CampaignID) (inviteDomain.Invite, error)
}
