package invite

import (
	campaignDomain "questmaster-core/internal/campaign/domain"
	inviteApp "questmaster-core/internal/invite/app"
	inviteDomain "questmaster-core/internal/invite/domain"
)

type GetInviteByCampaignIDUseCase struct {
	r inviteApp.InviteRepository
}

func NewGetInviteByCampaignID(r inviteApp.InviteRepository) *GetInviteByCampaignIDUseCase {
	return &GetInviteByCampaignIDUseCase{
		r: r,
	}
}

func (uc *GetInviteByCampaignIDUseCase) GetByCampaignID(campaignID campaignDomain.CampaignID) (*inviteDomain.Invite, error) {
	invite, err := uc.r.FindByCampaignID(campaignID)
	if err != nil {
		return &inviteDomain.Invite{}, err
	}

	return invite, nil
}
