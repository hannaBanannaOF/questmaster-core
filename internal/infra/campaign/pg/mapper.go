package campaign

import (
	domain "questmaster-core/internal/domain/campaign"
	"questmaster-core/internal/domain/rpg"
)

func MapRowToDomain(row CampaignRow) (domain.Campaign, error) {
	campaignName, err := domain.NewCampaignName(row.Name)
	if err != nil {
		return domain.Campaign{}, err
	}
	campaignSlug, err := rpg.NewSlug(row.Slug)
	if err != nil {
		return domain.Campaign{}, err
	}
	var campaignOverview *domain.CampaignOverview

	if row.Overview != nil {
		o, err := domain.NewCampaignOverview(*row.Overview)
		if err != nil {
			return domain.Campaign{}, err
		}
		campaignOverview = &o
	}
	return domain.Campaign{
		Id:       domain.CampaignID(row.Id),
		Name:     campaignName,
		Dm:       rpg.NewUserID(row.DmId),
		Status:   domain.CampaignStatus(row.Status),
		System:   rpg.System(row.System),
		Slug:     campaignSlug,
		Overview: campaignOverview,
	}, nil
}
