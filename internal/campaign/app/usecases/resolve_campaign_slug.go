package campaign

import (
	campaignApp "questmaster-core/internal/campaign/app"
	rpgApp "questmaster-core/internal/rpg/app"
)

type ResolveCampaignSlugUseCase struct {
	r campaignApp.CampaignRepository
}

func NewResolveCampaignSlug(r campaignApp.CampaignRepository) *ResolveCampaignSlugUseCase {
	return &ResolveCampaignSlugUseCase{r: r}
}

func (uc *ResolveCampaignSlugUseCase) Execute(cmd rpgApp.ResolveSlugCommand) (campaignApp.ResolveCampaignSlugReadModel, error) {
	campaign, err := uc.r.FindBySlug(cmd.Slug)
	if err != nil {
		return campaignApp.ResolveCampaignSlugReadModel{}, err
	}
	if campaign == nil {
		return campaignApp.ResolveCampaignSlugReadModel{}, ErrCampaignNotFound
	}

	return campaignApp.MapDomainToResolveCampaignSlugReadModel(*campaign), nil
}
