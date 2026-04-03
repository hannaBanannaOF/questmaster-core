package campaign

import (
	campaignApp "questmaster-core/internal/campaign/app"
	campaignDomain "questmaster-core/internal/campaign/domain"
)

type GetCampaignFromIDUseCase struct {
	r campaignApp.CampaignRepository
}

func NewGetCampaignFromID(r campaignApp.CampaignRepository) *GetCampaignFromIDUseCase {
	return &GetCampaignFromIDUseCase{r: r}
}

func (uc *GetCampaignFromIDUseCase) FindByID(campaignID campaignDomain.CampaignID) (campaignDomain.Campaign, error) {
	campaign, err := uc.r.FindById(campaignID)
	if err != nil {
		return campaignDomain.Campaign{}, err
	}
	if campaign == nil {
		return campaignDomain.Campaign{}, ErrCampaignNotFound
	}

	return *campaign, nil
}
