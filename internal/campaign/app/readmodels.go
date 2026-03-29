package campaign

import "github.com/google/uuid"

type CampaignListReadModel struct {
	Slug   string
	Name   string
	System string
	IsDM   bool
	Status string
}

type CampaignDetailsReadModel struct {
	Id         int
	Name       string
	Status     string
	System     string
	Slug       string
	Overview   *string
	IsDM       bool
	Characters []CampaignCharacterReadModel
}

type CampaignCharacterReadModel struct {
	Id   int
	Name string
}

type CreateCampaignReadModel struct {
	Slug string
}

type GetOrCreateInviteReadModel struct {
	InviteHash uuid.UUID
}

type ResolveCampaignSlugReadModel struct {
	ID int
}

type UpdateCampaignStatusReadModel struct {
	Status string
}
