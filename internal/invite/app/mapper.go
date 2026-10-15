package invite

import (
	campaignDomain "questmaster-core/internal/campaign/domain"
	characterDomain "questmaster-core/internal/character/domain"
	inviteDomain "questmaster-core/internal/invite/domain"
	userDomain "questmaster-core/internal/user/domain"
)

func MapDomainToInviteDetailsReadModel(inviteHash inviteDomain.InviteHash, campaign campaignDomain.Campaign, userID userDomain.UserID, premades []characterDomain.Character) InviteDetailReadModel {
	premadeModels := make([]InvitePremadeCharacterReadModel, 0, len(premades))
	for _, c := range premades {
		premadeModels = append(premadeModels, InvitePremadeCharacterReadModel{
			Slug: c.Slug.Value(),
			Name: c.Name.Value(),
		})
	}

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
		IsDM:                campaign.IsDM(userID),
		PremadeCharacters:   premadeModels,
	}
}
