package campaign

import (
	campaignDomain "questmaster-core/internal/campaign/domain"
	"questmaster-core/internal/shared/pagination"
	rpgDomain "questmaster-core/internal/rpg/domain"
	userDomain "questmaster-core/internal/user/domain"
)

type CampaignRepository interface {
	ListForUser(userID userDomain.UserID, filters campaignDomain.CampaignListFilters, page pagination.Page) (pagination.Result[campaignDomain.Campaign], error)
	CountByStatusForUser(userID userDomain.UserID, filters campaignDomain.CampaignListFilters) (map[campaignDomain.CampaignStatus]int, error)
	FindBySlug(slug rpgDomain.Slug) (*campaignDomain.Campaign, error)
	FindById(id campaignDomain.CampaignID) (*campaignDomain.Campaign, error)
	Create(Name campaignDomain.CampaignName, Overview *campaignDomain.CampaignOverview, DmID userDomain.UserID, System rpgDomain.System) (campaignDomain.Campaign, error)
	UpdateStatus(newStatus campaignDomain.CampaignStatus, id campaignDomain.CampaignID) (campaignDomain.Campaign, error)
	DeleteById(id campaignDomain.CampaignID) (bool, error)
}
