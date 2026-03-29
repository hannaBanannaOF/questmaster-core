package campaign

import (
	campaignApp "questmaster-core/internal/campaign/app"
)

type DeleteCampaignUseCase struct {
	r campaignApp.CampaignRepository
}

func NewDeleteCampaign(r campaignApp.CampaignRepository) *DeleteCampaignUseCase {
	return &DeleteCampaignUseCase{
		r: r,
	}
}

func (uc *DeleteCampaignUseCase) Execute(cmd campaignApp.DeleteCampaignCommand) error {
	campaign, err := uc.r.FindById(cmd.ID)
	if err != nil {
		return err
	}

	if campaign == nil {
		return ErrCampaignNotFound
	}

	if err := campaign.CanDelete(cmd.UserID); err != nil {
		return err
	}

	deleted, err := uc.r.DeleteById(cmd.ID)
	if err != nil {
		return err
	}
	if !deleted {
		return ErrCampaignNotFound
	}

	return nil
}
