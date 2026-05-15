package campaign

import (
	"errors"
	campaignApp "questmaster-core/internal/campaign/app"
	inviteUsecases "questmaster-core/internal/invite/app/usecases"
)

type DeleteCampaignUseCase struct {
	r              campaignApp.CampaignRepository
	deleteInviteUc CampaignInviteDeleter
}

func NewDeleteCampaign(r campaignApp.CampaignRepository, deleteInviteUc CampaignInviteDeleter) *DeleteCampaignUseCase {
	return &DeleteCampaignUseCase{
		r:              r,
		deleteInviteUc: deleteInviteUc,
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

	if err := uc.deleteInviteUc.DeleteByCampaignID(cmd.ID); err != nil && !errors.Is(err, inviteUsecases.ErrInviteNotFound) {
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
