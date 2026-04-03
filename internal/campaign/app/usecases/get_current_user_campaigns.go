package campaign

import (
	campaignApp "questmaster-core/internal/campaign/app"
	campaigndomain "questmaster-core/internal/campaign/domain"
)

type GetCurrentUserCampaignsUseCase struct {
	r campaignApp.CampaignRepository
}

func NewGetCurrentUserMyCampaigns(r campaignApp.CampaignRepository) *GetCurrentUserCampaignsUseCase {
	return &GetCurrentUserCampaignsUseCase{r: r}
}

func (uc *GetCurrentUserCampaignsUseCase) Execute(cmd campaignApp.GetCurrentUserCampaignsCommand) ([]campaigndomain.Campaign, error) {
	dmCampaigns, err := uc.r.GetByDmId(cmd.UserID)
	if err != nil {
		return nil, err
	}

	playerCampaigns, err := uc.r.GetByPlayerId(cmd.UserID)
	if err != nil {
		return nil, err
	}

	seen := make(map[campaigndomain.CampaignID]struct{})
	items := make([]campaigndomain.Campaign, 0)

	for _, c := range dmCampaigns {
		seen[c.Id] = struct{}{}
		items = append(items, c)
	}

	for _, c := range playerCampaigns {
		if _, exists := seen[c.Id]; exists {
			continue
		}
		items = append(items, c)
	}

	return items, nil
}
