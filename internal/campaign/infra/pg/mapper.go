package campaign

import (
	campaignDomain "questmaster-core/internal/campaign/domain"
	rpgDomain "questmaster-core/internal/rpg/domain"
	userDomain "questmaster-core/internal/user/domain"
)

func MapRowToDomain(row CampaignRow) (campaignDomain.Campaign, error) {
	name, err := campaignDomain.NewCampaignName(row.Name)
	if err != nil {
		return campaignDomain.Campaign{}, err
	}
	slug, err := rpgDomain.NewSlug(row.Slug)
	if err != nil {
		return campaignDomain.Campaign{}, err
	}
	var overview *campaignDomain.CampaignOverview

	if row.Overview != nil {
		o := campaignDomain.NewCampaignOverview(*row.Overview)
		overview = &o
	}

	status, err := campaignDomain.NewCampaignStatus(row.Status)
	if err != nil {
		return campaignDomain.Campaign{}, err
	}

	system, err := rpgDomain.NewSystem(row.System)
	if err != nil {
		return campaignDomain.Campaign{}, err
	}

	return campaignDomain.Campaign{
		Id:          campaignDomain.NewCampaignID(row.ID),
		Name:        name,
		Dm:          userDomain.NewUserID(row.DmID),
		Status:      status,
		System:      system,
		Slug:        slug,
		Overview:    overview,
		PlayerCount: campaignDomain.NewPlayerCount(row.PlayerCount),
	}, nil
}
