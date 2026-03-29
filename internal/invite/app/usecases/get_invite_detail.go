package invite

import (
	inviteApp "questmaster-core/internal/invite/app"
)

type GetInviteDetailUseCase struct {
	r                        inviteApp.InviteRepository
	getCampaignUC            InviteCampaignFinder
	getAvailableCharactersUC InviteAvailableCharacterFinder
}

func NewGetInviteDetail(
	r inviteApp.InviteRepository,
	getCampaignUC InviteCampaignFinder,
	getAvailableCharactersUC InviteAvailableCharacterFinder,
) *GetInviteDetailUseCase {
	return &GetInviteDetailUseCase{
		r:                        r,
		getCampaignUC:            getCampaignUC,
		getAvailableCharactersUC: getAvailableCharactersUC,
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

	availableCharacters, err := uc.getAvailableCharactersUC.GetBySystemAndCampaignIDNull(cmd.UserID, campaign.System)
	if err != nil {
		return inviteApp.InviteDetailReadModel{}, err
	}

	return inviteApp.MapDomainToInviteDetailsReadModel(campaign, availableCharacters), nil
}
