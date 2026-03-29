package campaign

import (
	campaignDomain "questmaster-core/internal/campaign/domain"
	characterDomain "questmaster-core/internal/character/domain"
	inviteDomain "questmaster-core/internal/invite/domain"
	rpgDomain "questmaster-core/internal/rpg/domain"
)

func MapDomainToListReadModel(domain campaignDomain.Campaign, userId rpgDomain.UserID) CampaignListReadModel {
	return CampaignListReadModel{
		Slug:   domain.Slug.Value(),
		Name:   domain.Name.Value(),
		System: domain.System.Value(),
		IsDM:   domain.IsDM(userId),
		Status: domain.Status.Value(),
	}
}

func MapDomainToDetailReadModel(campaign campaignDomain.Campaign, characters []characterDomain.Character, userId rpgDomain.UserID) CampaignDetailsReadModel {

	charactersModels := make([]CampaignCharacterReadModel, 0, len(characters))
	for _, c := range characters {
		charactersModels = append(charactersModels, CampaignCharacterReadModel{
			Id:   c.Id.Value(),
			Name: c.Name.Value(),
		})
	}

	var campaignOverview *string
	if campaign.Overview != nil {
		o := campaign.Overview.Value()
		campaignOverview = &o
	}
	return CampaignDetailsReadModel{
		Id:         campaign.Id.Value(),
		Name:       campaign.Name.Value(),
		Status:     campaign.Status.Value(),
		System:     campaign.System.Value(),
		Slug:       campaign.Slug.Value(),
		Overview:   campaignOverview,
		IsDM:       campaign.IsDM(userId),
		Characters: charactersModels,
	}
}

func MapDomainToCreateReadModel(campaign campaignDomain.Campaign) CreateCampaignReadModel {
	return CreateCampaignReadModel{
		Slug: campaign.Slug.Value(),
	}
}

func MapDomainToGetOrCreateInviteReadModel(invite inviteDomain.Invite) GetOrCreateInviteReadModel {
	return GetOrCreateInviteReadModel{
		InviteHash: invite.Hash.Value(),
	}
}

func MapDomainToResolveCampaignSlugReadModel(campaign campaignDomain.Campaign) ResolveCampaignSlugReadModel {
	return ResolveCampaignSlugReadModel{
		ID: campaign.Id.Value(),
	}
}

func MapDomainToUpdateCampaignStatusReadModel(campaign campaignDomain.Campaign) UpdateCampaignStatusReadModel {
	return UpdateCampaignStatusReadModel{
		Status: campaign.Status.Value(),
	}
}
