package campaign

import (
	campaignApp "questmaster-core/internal/campaign/app"
	rpgDomain "questmaster-core/internal/rpg/domain"
)

type ResolveCampaignSlugUseCase struct {
	r campaignApp.CampaignRepository
}

func NewResolveCampaignSlug(r campaignApp.CampaignRepository) *ResolveCampaignSlugUseCase {
	return &ResolveCampaignSlugUseCase{r: r}
}

func (uc *ResolveCampaignSlugUseCase) Execute(slug rpgDomain.Slug) (campaignApp.ResolveCampaignSlugReadModel, error) {
	campaign, err := uc.r.FindBySlug(slug)
	if err != nil {
		return campaignApp.ResolveCampaignSlugReadModel{}, err
	}
	if campaign == nil {
		return campaignApp.ResolveCampaignSlugReadModel{}, ErrCampaignNotFound
	}

	return campaignApp.MapDomainToResolveCampaignSlugReadModel(*campaign), nil
}
