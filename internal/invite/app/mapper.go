package invite

import (
	campaignDomain "questmaster-core/internal/campaign/domain"
	inviteDomain "questmaster-core/internal/invite/domain"
)

func MapDomainToInviteDetailsReadModel(inviteHash inviteDomain.InviteHash, campaign campaignDomain.Campaign) InviteDetailReadModel {
	var overview *string
	if campaign.Overview != nil {
		o := campaign.Overview.Value()
		overview = &o
	}

	return InviteDetailReadModel{
		InviteHash:          inviteHash.Value(),
		CampaignSlug:        campaign.Slug.Value(),
		CampaignName:        campaign.Name.Value(),
		CampaignSystem:      campaign.System.Value(),
		CampaignOverview:    overview,
		CampaignPlayerCount: campaign.PlayerCount.Value(),
	}
}
