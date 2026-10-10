package campaign

import (
	campaignDomain "questmaster-core/internal/campaign/domain"
	characterDomain "questmaster-core/internal/character/domain"
	inviteDomain "questmaster-core/internal/invite/domain"
	userDomain "questmaster-core/internal/user/domain"
)

type CampaignCharacterFinder interface {
	GetByCampaignID(campaignID campaignDomain.CampaignID) ([]characterDomain.Character, error)
}

type CampaignPlayerCharactersFinder interface {
	GetByPlayerInCampaigns(userID userDomain.UserID, campaignIDs []campaignDomain.CampaignID) ([]characterDomain.Character, error)
}

type CampaignInviteFinder interface {
	GetByCampaignID(campaignID campaignDomain.CampaignID) (*inviteDomain.Invite, error)
}
