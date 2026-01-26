package campaign

import (
	domain "questmaster-core/internal/domain/campaign"
	"questmaster-core/internal/domain/rpg"
)

type CampaignRepository interface {
	GetByDmId(userId rpg.UserID) ([]domain.Campaign, error)
	GetByPlayerId(userId rpg.UserID) ([]domain.Campaign, error)
	FindBySlug(slug rpg.Slug) (*domain.Campaign, error)
	FindById(id domain.CampaignID) (*domain.Campaign, error)
	Create(input CreateCampaignInput) (domain.Campaign, error)
	UpdateStatus(newStatus domain.CampaignStatus, id domain.CampaignID) (domain.Campaign, error)
}
