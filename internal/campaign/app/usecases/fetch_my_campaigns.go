package campaign

import (
	campaignApp "questmaster-core/internal/campaign/app"
	campaigndomain "questmaster-core/internal/campaign/domain"
	rpgDomain "questmaster-core/internal/rpg/domain"
)

type FetchMyCampaignsUseCase struct {
	r campaignApp.CampaignRepository
}

func NewFetchMyCampaigns(r campaignApp.CampaignRepository) *FetchMyCampaignsUseCase {
	return &FetchMyCampaignsUseCase{r: r}
}

func (uc *FetchMyCampaignsUseCase) Execute(userID rpgDomain.UserID) ([]campaignApp.CampaignListReadModel, error) {
	dmCampaigns, err := uc.r.GetByDmId(userID)
	if err != nil {
		return nil, err
	}

	playerCampaigns, err := uc.r.GetByPlayerId(userID)
	if err != nil {
		return nil, err
	}

	seen := make(map[campaigndomain.CampaignID]struct{})
	items := make([]campaignApp.CampaignListReadModel, 0)

	for _, c := range dmCampaigns {
		seen[c.Id] = struct{}{}
		items = append(items, campaignApp.MapDomainToListReadModel(c, userID))
	}

	for _, c := range playerCampaigns {
		if _, exists := seen[c.Id]; exists {
			continue
		}
		items = append(items, campaignApp.MapDomainToListReadModel(c, userID))
	}

	return items, nil
}
