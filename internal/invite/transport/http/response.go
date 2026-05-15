package invite

import "github.com/google/uuid"

type InviteDetailsResponse struct {
	InviteHash          uuid.UUID `json:"invite_hash"`
	CampaignSlug        string    `json:"campaign_slug"`
	CampaignName        string    `json:"campaign_name"`
	CampaignOverview    *string   `json:"campaign_overview"`
	CampaignSystem      string    `json:"campaign_system"`
	CampaignPlayerCount int       `json:"campaign_player_count"`
}

type InviteCreateResponse struct {
	Hash uuid.UUID `json:"hash"`
}
