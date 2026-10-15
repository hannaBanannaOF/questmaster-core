package campaign

import (
	campaignDomain "questmaster-core/internal/campaign/domain"
	rpgDomain "questmaster-core/internal/rpg/domain"
	"questmaster-core/internal/shared/pagination"
	userDomain "questmaster-core/internal/user/domain"
)

type CreateCampaignCommand struct {
	Name     campaignDomain.CampaignName
	Overview *campaignDomain.CampaignOverview
	System   rpgDomain.System
	DmID     userDomain.UserID
	DmName   rpgDomain.TableName
}

type UpdateCampaignDmNameCommand struct {
	CampaignID campaignDomain.CampaignID
	UserID     userDomain.UserID
	DmName     rpgDomain.TableName
}

type DeleteCampaignCommand struct {
	ID     campaignDomain.CampaignID
	UserID userDomain.UserID
}

type GetCurrentUserCampaignsCommand struct {
	UserID  userDomain.UserID
	Filters campaignDomain.CampaignListFilters
	Page    pagination.Page
}

type GetCampaignStatusCountsCommand struct {
	UserID  userDomain.UserID
	Filters campaignDomain.CampaignListFilters
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
