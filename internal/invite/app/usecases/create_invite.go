package invite

import (
	inviteApp "questmaster-core/internal/invite/app"
	inviteDomain "questmaster-core/internal/invite/domain"
)

type CreateInviteUseCase struct {
	r             inviteApp.InviteRepository
	getCampaignUC InviteCampaignFinder
}

func NewCreateInvite(r inviteApp.InviteRepository, getCampaignUC InviteCampaignFinder) *CreateInviteUseCase {
	return &CreateInviteUseCase{
		r:             r,
		getCampaignUC: getCampaignUC,
	}
}

func (uc *CreateInviteUseCase) Execute(cmd inviteApp.CreateInviteCommand) (inviteDomain.Invite, error) {
	campaign, err := uc.getCampaignUC.FindByID(cmd.CampaignID)
	if err != nil {
		return inviteDomain.Invite{}, err
	}

	if err := campaign.CanInvite(cmd.UserID); err != nil {
		return inviteDomain.Invite{}, err
	}

	invite, err := uc.r.Create(cmd.CampaignID)
	if err != nil {
		return inviteDomain.Invite{}, err
	}

	if invite == nil {
		return inviteDomain.Invite{}, ErrInviteAlreadyExists
	}

	return *invite, nil
}
