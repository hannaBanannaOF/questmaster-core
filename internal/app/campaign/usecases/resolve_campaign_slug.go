package campaign

import (
	app "questmaster-core/internal/app/campaign"
	"questmaster-core/internal/domain/rpg"
)

type ResolveCampaignSlugUseCase struct {
	r app.CampaignRepository
}

func NewResolveCampaignSlug(r app.CampaignRepository) *ResolveCampaignSlugUseCase {
	return &ResolveCampaignSlugUseCase{r: r}
}

func (uc *ResolveCampaignSlugUseCase) Execute(slug string) (int, error) {
	slugDomain, err := rpg.NewSlug(slug)
	if err != nil {
		return 0, err
	}
	domain, err := uc.r.FindBySlug(slugDomain)
	if err != nil {
		return 0, err
	}
	if domain == nil {
		return 0, ErrCampaignNotFound
	}

	return int(domain.Id), nil
}
