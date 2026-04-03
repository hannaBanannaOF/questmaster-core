package campaign

import (
	campaignDomain "questmaster-core/internal/campaign/domain"
	rpgDomain "questmaster-core/internal/rpg/domain"
	userDomain "questmaster-core/internal/user/domain"
)

type CreateCampaignCommand struct {
	Name     campaignDomain.CampaignName
	Overview *campaignDomain.CampaignOverview
	System   rpgDomain.System
	DmID     userDomain.UserID
}

type DeleteCampaignCommand struct {
	ID     campaignDomain.CampaignID
	UserID userDomain.UserID
}

type GetCurrentUserCampaignsCommand struct {
	UserID userDomain.UserID
}

type GetCampaignDetailsCommand struct {
	ID     campaignDomain.CampaignID
	UserID userDomain.UserID
}

type UpdateCampaignStatusCommand struct {
	CampaignID campaignDomain.CampaignID
	UserID     userDomain.UserID
	NewStatus  campaignDomain.CampaignStatus
}
