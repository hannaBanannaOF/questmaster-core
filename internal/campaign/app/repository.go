package campaign

import (
	campaignDomain "questmaster-core/internal/campaign/domain"
	rpgDomain "questmaster-core/internal/rpg/domain"
	userDomain "questmaster-core/internal/user/domain"
)

type CampaignRepository interface {
	GetByDmId(userID userDomain.UserID) ([]campaignDomain.Campaign, error)
	GetByPlayerId(userID userDomain.UserID) ([]campaignDomain.Campaign, error)
	FindBySlug(slug rpgDomain.Slug) (*campaignDomain.Campaign, error)
	FindById(id campaignDomain.CampaignID) (*campaignDomain.Campaign, error)
	Create(Name campaignDomain.CampaignName, Overview *campaignDomain.CampaignOverview, DmID userDomain.UserID, System rpgDomain.System) (campaignDomain.Campaign, error)
	UpdateStatus(newStatus campaignDomain.CampaignStatus, id campaignDomain.CampaignID) (campaignDomain.Campaign, error)
	DeleteById(id campaignDomain.CampaignID) (bool, error)
}
