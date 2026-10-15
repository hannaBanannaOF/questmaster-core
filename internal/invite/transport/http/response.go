package invite

import "github.com/google/uuid"

type InviteDetailsResponse struct {
	InviteHash          uuid.UUID `json:"invite_hash"`
	CampaignSlug        string    `json:"campaign_slug"`
	CampaignName        string    `json:"campaign_name"`
	CampaignOverview    *string   `json:"campaign_overview"`
	CampaignSystem      string    `json:"campaign_system"`
	CampaignPlayerCount int       `json:"campaign_player_count"`
	IsDM                bool      `json:"is_dm"`
	// Unclaimed premade characters the requester may choose when accepting, ordered by name
	PremadeCharacters []InvitePremadeCharacterResponse `json:"premade_characters"`
}

type InvitePremadeCharacterResponse struct {
	Slug string `json:"slug"`
	Name string `json:"name"`
}

type InviteCreateResponse struct {
	Hash uuid.UUID `json:"hash"`
}
