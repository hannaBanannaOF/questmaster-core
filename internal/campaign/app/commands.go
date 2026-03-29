package campaign

import (
	campaignDomain "questmaster-core/internal/campaign/domain"
	rpgDomain "questmaster-core/internal/rpg/domain"
)

type CreateCampaignCommand struct {
	Name     campaignDomain.CampaignName
	Overview *campaignDomain.CampaignOverview
	System   rpgDomain.System
	DmID     rpgDomain.UserID
}

type DeleteCampaignCommand struct {
	ID     campaignDomain.CampaignID
	UserID rpgDomain.UserID
}

type GetCampaignDetailsCommand struct {
	ID     campaignDomain.CampaignID
	UserID rpgDomain.UserID
}

type GetOrCreateCampaignInviteCommand struct {
	CampaignID campaignDomain.CampaignID
	UserID     rpgDomain.UserID
}

type UpdateCampaignStatusCommand struct {
	CampaignID campaignDomain.CampaignID
	UserID     rpgDomain.UserID
	NewStatus  campaignDomain.CampaignStatus
}
