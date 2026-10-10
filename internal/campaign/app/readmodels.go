package campaign

import "github.com/google/uuid"

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
