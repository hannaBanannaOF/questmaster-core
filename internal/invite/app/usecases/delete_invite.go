package invite

import (
	campaignDomain "questmaster-core/internal/campaign/domain"
	inviteApp "questmaster-core/internal/invite/app"
)

type DeleteInviteUseCase struct {
	r inviteApp.InviteRepository
}

func NewDeleteInvite(r inviteApp.InviteRepository) *DeleteInviteUseCase {
	return &DeleteInviteUseCase{
		r: r,
	}
}

func (uc *DeleteInviteUseCase) DeleteByCampaignID(campaignID campaignDomain.CampaignID) error {
	deleted, err := uc.r.DeleteByCampaignID(campaignID)
	if err != nil {
		return err
	}
	if !deleted {
		return ErrInviteNotFound
	}

	return nil
}
