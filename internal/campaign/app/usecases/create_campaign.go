package campaign

import (
	campaignApp "questmaster-core/internal/campaign/app"
)

type CreateCampaignUseCase struct {
	r campaignApp.CampaignRepository
}

func NewCreateCampaign(r campaignApp.CampaignRepository) *CreateCampaignUseCase {
	return &CreateCampaignUseCase{r: r}
}

func (uc *CreateCampaignUseCase) Execute(cmd campaignApp.CreateCampaignCommand) (campaignApp.CreateCampaignReadModel, error) {
	campaign, err := uc.r.Create(cmd.Name, cmd.Overview, cmd.DmID, cmd.System)
	if err != nil {
		return campaignApp.CreateCampaignReadModel{}, err
	}

	return campaignApp.MapDomainToCreateReadModel(campaign), nil
}
