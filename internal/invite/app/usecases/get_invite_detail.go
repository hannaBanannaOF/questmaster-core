package invite

import (
	inviteApp "questmaster-core/internal/invite/app"
)

type GetInviteDetailUseCase struct {
	r             inviteApp.InviteRepository
	getCampaignUC InviteCampaignFinder
}

func NewGetInviteDetail(
	r inviteApp.InviteRepository,
	getCampaignUC InviteCampaignFinder,
) *GetInviteDetailUseCase {
	return &GetInviteDetailUseCase{
		r:             r,
		getCampaignUC: getCampaignUC,
	}
}

func (uc *GetInviteDetailUseCase) Execute(cmd inviteApp.GetinviteDetailsCommand) (inviteApp.InviteDetailReadModel, error) {
	invite, err := uc.r.FindByHash(cmd.Hash)
	if err != nil {
		return inviteApp.InviteDetailReadModel{}, err
	}
	if invite == nil {
		return inviteApp.InviteDetailReadModel{}, ErrInviteNotFound
	}

	campaign, err := uc.getCampaignUC.FindByID(invite.CampaignId)
	if err != nil {
		return inviteApp.InviteDetailReadModel{}, err
	}

	return inviteApp.MapDomainToInviteDetailsReadModel(cmd.Hash, campaign), nil
}
