package invite

import (
	campaignDomain "questmaster-core/internal/campaign/domain"
	inviteDomain "questmaster-core/internal/invite/domain"
)

func MapRowToDomain(row InviteRow) (inviteDomain.Invite, error) {
	return inviteDomain.Invite{
		Id:         inviteDomain.InviteID(row.Id),
		CampaignId: campaignDomain.CampaignID(row.CampaignId),
		Hash:       inviteDomain.NewHash(row.Hash),
	}, nil
}
