package invite

import (
	campaignDomain "questmaster-core/internal/campaign/domain"
	characterDomain "questmaster-core/internal/character/domain"
)

func MapDomainToInviteDetailsReadModel(campaign campaignDomain.Campaign, characters []characterDomain.Character) InviteDetailReadModel {
	charactersRm := make([]InviteDetailCharacterItem, 0, len(characters))

	for _, c := range characters {
		charactersRm = append(charactersRm, InviteDetailCharacterItem{
			ID:   c.Id.Value(),
			Name: c.Name.Value(),
		})
	}

	return InviteDetailReadModel{
		CampaignID:   campaign.Id.Value(),
		CampaignName: campaign.Name.Value(),
		Characters:   charactersRm,
	}
}
