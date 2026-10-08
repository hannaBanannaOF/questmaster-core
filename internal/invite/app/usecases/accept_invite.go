package invite

import (
	inviteApp "questmaster-core/internal/invite/app"
)

type AcceptInviteUseCase struct {
	r                         inviteApp.InviteRepository
	getCampaignUC             InviteCampaignFinder
	linkCharacterToCampaignUC InviteCharacterCampaignLinker
}

func NewAcceptInvite(
	r inviteApp.InviteRepository,
	getCampaignUC InviteCampaignFinder,
	linkCharacterToCampaignUC InviteCharacterCampaignLinker,
) *AcceptInviteUseCase {
	return &AcceptInviteUseCase{
		r:                         r,
		getCampaignUC:             getCampaignUC,
		linkCharacterToCampaignUC: linkCharacterToCampaignUC,
	}
}

func (uc *AcceptInviteUseCase) Execute(cmd inviteApp.AcceptInviteCommand) error {
	invite, err := uc.r.FindByHash(cmd.Hash)
	if err != nil {
		return err
	}

	if invite == nil {
		return ErrInviteNotFound
	}

	campaign, err := uc.getCampaignUC.FindByID(invite.CampaignId)
	if err != nil {
		return err
	}

	if err := campaign.CanJoin(cmd.UserID); err != nil {
		return err
	}

	_, err = uc.linkCharacterToCampaignUC.LinkToCampaign(invite.CampaignId, cmd.CharacterSlug, cmd.UserID)
	if err != nil {
		return err
	}

	return nil
}
