package invite

import (
	inviteApp "questmaster-core/internal/invite/app"
)

type AcceptInviteUseCase struct {
	r                         inviteApp.InviteRepository
	linkCharacterToCampaignUC InviteCharacterCampaignLinker
}

func NewAcceptInvite(r inviteApp.InviteRepository, linkCharacterToCampaignUC InviteCharacterCampaignLinker) *AcceptInviteUseCase {
	return &AcceptInviteUseCase{
		r:                         r,
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

	_, err = uc.linkCharacterToCampaignUC.LinkToCampaign(invite.CampaignId, cmd.CharacterSheetID)
	if err != nil {
		return err
	}

	return nil
}
