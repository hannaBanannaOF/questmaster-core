package invite

import "github.com/google/uuid"

type InviteDetailsCharacterListItem struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type InviteDetailsResponse struct {
	CampaignID          int     `json:"campaign_id"`
	CampaignName        string  `json:"campaign_name"`
	CampaignOverview    *string `json:"campaign_overview"`
	CampaignSystem      string  `json:"campaign_system"`
	CampaignPlayerCount int     `json:"campaign_player_count"`
}

type InviteCreateResponse struct {
	Hash uuid.UUID `json:"hash"`
}
