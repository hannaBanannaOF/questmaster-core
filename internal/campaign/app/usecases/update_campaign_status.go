package campaign

import (
	campaignApp "questmaster-core/internal/campaign/app"
)

type UpdateCampaignStatusUseCase struct {
	r campaignApp.CampaignRepository
}

func NewUpdateStatus(r campaignApp.CampaignRepository) *UpdateCampaignStatusUseCase {
	return &UpdateCampaignStatusUseCase{r: r}
}

func (uc *UpdateCampaignStatusUseCase) Execute(cmd campaignApp.UpdateCampaignStatusCommand) (campaignApp.UpdateCampaignStatusReadModel, error) {
	campaign, err := uc.r.FindById(cmd.CampaignID)
	if err != nil {
		return campaignApp.UpdateCampaignStatusReadModel{}, err
	}
	if campaign == nil {
		return campaignApp.UpdateCampaignStatusReadModel{}, ErrCampaignNotFound
	}

	if err := campaign.ChangeStatus(cmd.NewStatus, cmd.UserID); err != nil {
		return campaignApp.UpdateCampaignStatusReadModel{}, err
	}

	newCampaign, err := uc.r.UpdateStatus(campaign.Status, campaign.Id)
	if err != nil {
		return campaignApp.UpdateCampaignStatusReadModel{}, err
	}

	return campaignApp.MapDomainToUpdateCampaignStatusReadModel(newCampaign), nil
}
