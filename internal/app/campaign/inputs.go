package campaign

import (
	domain "questmaster-core/internal/domain/campaign"
	"questmaster-core/internal/domain/rpg"
)

type CreateCampaignInput struct {
	Name     domain.CampaignName
	Overview *domain.CampaignOverview
	System   rpg.System
	Dm       rpg.UserID
}
