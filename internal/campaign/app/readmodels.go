package campaign

import (
	campaignDomain "questmaster-core/internal/campaign/domain"

	"github.com/google/uuid"
)

type CampaignDetailsReadModel struct {
	Id         int
	Name       string
	Status     string
	System     string
	Slug       string
	Overview   *string
	IsDM       bool
	Characters []CampaignCharacterReadModel
	InviteHash *uuid.UUID
}

type CampaignCharacterReadModel struct {
	Id        int
	Name      string
	CurrentHP *int
}

// CampaignListItemReadModel is a campaign of the requester list, with the requester characters in it
type CampaignListItemReadModel struct {
	Campaign     campaignDomain.Campaign
	MyCharacters []CampaignCharacterRefReadModel
}

type CampaignCharacterRefReadModel struct {
	Slug string
	Name string
}

type CreateCampaignReadModel struct {
	Slug string
}

type ResolveCampaignSlugReadModel struct {
	ID int
}

// CampaignStatusCountsReadModel has an entry for every campaign status
type CampaignStatusCountsReadModel map[string]int

type UpdateCampaignStatusReadModel struct {
	Status string
}
