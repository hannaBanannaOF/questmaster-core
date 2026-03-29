package invite

import (
	campaignDomain "questmaster-core/internal/campaign/domain"
	inviteApp "questmaster-core/internal/invite/app"
	inviteDomain "questmaster-core/internal/invite/domain"
)

type CreateInviteUseCase struct {
	r inviteApp.InviteRepository
}

func NewCreateInvite(r inviteApp.InviteRepository) *CreateInviteUseCase {
	return &CreateInviteUseCase{
		r: r,
	}
}

func (uc *CreateInviteUseCase) Create(campaignID campaignDomain.CampaignID) (inviteDomain.Invite, error) {
	invite, err := uc.r.Create(campaignID)
	if err != nil {
		return inviteDomain.Invite{}, err
	}

	if invite == nil {
		return inviteDomain.Invite{}, ErrInviteAlreadyExists
	}

	return *invite, nil
}
