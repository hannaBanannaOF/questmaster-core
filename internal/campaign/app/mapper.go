package campaign

import (
	campaignDomain "questmaster-core/internal/campaign/domain"
	characterDomain "questmaster-core/internal/character/domain"
	inviteDomain "questmaster-core/internal/invite/domain"
	userDomain "questmaster-core/internal/user/domain"

	"github.com/google/uuid"
)

type CampaignDetailsInput struct {
	Campaign   campaignDomain.Campaign
	Characters []characterDomain.Character
	// Unclaimed premade characters; only loaded for the DM
	PremadeCharacters []characterDomain.Character
	Invite            *inviteDomain.Invite
	UserID            userDomain.UserID
}

func mapCampaignCharacters(characters []characterDomain.Character, userID userDomain.UserID) []CampaignCharacterReadModel {
	models := make([]CampaignCharacterReadModel, 0, len(characters))
	for _, c := range characters {
		models = append(models, CampaignCharacterReadModel{
			Id:     c.Id.Value(),
			Slug:   c.Slug.Value(),
			Name:   c.Name.Value(),
			IsMine: c.IsPlayer(userID),
		})
	}
	return models
}

func MapDomainToDetailReadModel(input CampaignDetailsInput) CampaignDetailsReadModel {
	charactersModels := mapCampaignCharacters(input.Characters, input.UserID)

	var campaignOverview *string
	if input.Campaign.Overview != nil {
		o := input.Campaign.Overview.Value()
		campaignOverview = &o
	}

	var inviteHash *uuid.UUID
	if input.Invite != nil {
		h := input.Invite.Hash.Value()
		inviteHash = &h
	}

	return CampaignDetailsReadModel{
		Id:                input.Campaign.Id.Value(),
		Name:              input.Campaign.Name.Value(),
		Status:            input.Campaign.Status.Value(),
		System:            input.Campaign.System.Value(),
		Slug:              input.Campaign.Slug.Value(),
		Overview:          campaignOverview,
		IsDM:              input.Campaign.IsDM(input.UserID),
		Characters:        charactersModels,
		PremadeCharacters: mapCampaignCharacters(input.PremadeCharacters, input.UserID),
		InviteHash:        inviteHash,
	}
}

func MapDomainToCreateReadModel(campaign campaignDomain.Campaign) CreateCampaignReadModel {
	return CreateCampaignReadModel{
		Slug: campaign.Slug.Value(),
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
