package campaign

import (
	app "questmaster-core/internal/app/campaign"
	domain "questmaster-core/internal/domain/campaign"
	"questmaster-core/internal/domain/rpg"

	"github.com/google/uuid"
)

type FetchMyCampaignsUseCase struct {
	r app.CampaignRepository
}

func NewFetchMyCampaigns(r app.CampaignRepository) *FetchMyCampaignsUseCase {
	return &FetchMyCampaignsUseCase{r: r}
}

func (uc *FetchMyCampaignsUseCase) Execute(userId uuid.UUID) ([]app.CampaignListReadModel, error) {
	dmCampaigns, err := uc.r.GetByDmId(rpg.NewUserID(userId))
	if err != nil {
		return nil, err
	}

	playerCampaigns, err := uc.r.GetByPlayerId(rpg.NewUserID(userId))
	if err != nil {
		return nil, err
	}

	seen := make(map[domain.CampaignID]struct{})
	items := make([]app.CampaignListReadModel, 0)

	for _, c := range dmCampaigns {
		seen[c.Id] = struct{}{}
		items = append(items, app.MapDomainToListReadModel(c, userId))
	}

	for _, c := range playerCampaigns {
		if _, exists := seen[c.Id]; exists {
			continue
		}
		items = append(items, app.MapDomainToListReadModel(c, userId))
	}

	return items, nil
}
