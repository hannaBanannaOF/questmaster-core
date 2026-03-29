package campaign

import (
	"errors"
	campaignApp "questmaster-core/internal/campaign/app"
	inviteUsecases "questmaster-core/internal/invite/app/usecases"
)

type GetOrCreateCampaignInviteUseCase struct {
	getCampaignByIdUc       GetCampaignFromIDUseCase
	getInviteByCampaignIdUc CampaignInviteFinder
	createInviteUc          CampaignInviteCreator
}

func NewGetOrCreateCampaignInvite(
	getCampaignByIdUc GetCampaignFromIDUseCase,
	getInviteByCampaignIdUc CampaignInviteFinder,
	createInviteUc CampaignInviteCreator,
) *GetOrCreateCampaignInviteUseCase {
	return &GetOrCreateCampaignInviteUseCase{
		getCampaignByIdUc:       getCampaignByIdUc,
		getInviteByCampaignIdUc: getInviteByCampaignIdUc,
		createInviteUc:          createInviteUc,
	}
}

func (uc *GetOrCreateCampaignInviteUseCase) Execute(cmd campaignApp.GetOrCreateCampaignInviteCommand) (campaignApp.GetOrCreateInviteReadModel, error) {
	campaign, err := uc.getCampaignByIdUc.FindByID(cmd.CampaignID)
	if err != nil {
		return campaignApp.GetOrCreateInviteReadModel{}, err
	}

	if err := campaign.CanEdit(cmd.UserID); err != nil {
		invite, getErr := uc.getInviteByCampaignIdUc.GetByCampaignID(cmd.CampaignID)
		if getErr != nil {
			return campaignApp.GetOrCreateInviteReadModel{}, getErr
		}
		return campaignApp.MapDomainToGetOrCreateInviteReadModel(invite), nil
	}

	invite, err := uc.createInviteUc.Create(cmd.CampaignID)

	if err == nil {
		return campaignApp.MapDomainToGetOrCreateInviteReadModel(invite), nil
	}

	if errors.Is(err, inviteUsecases.ErrInviteAlreadyExists) {

		invite, getErr := uc.getInviteByCampaignIdUc.GetByCampaignID(cmd.CampaignID)
		if getErr != nil {
			return campaignApp.GetOrCreateInviteReadModel{}, getErr
		}

		return campaignApp.MapDomainToGetOrCreateInviteReadModel(invite), inviteUsecases.ErrInviteAlreadyExists
	}

	return campaignApp.GetOrCreateInviteReadModel{}, err
}
