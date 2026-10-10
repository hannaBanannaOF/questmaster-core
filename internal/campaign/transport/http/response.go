package campaign

import "github.com/google/uuid"

type CampaignListResponse struct {
	Slug        string `json:"slug"`
	Name        string `json:"name"`
	IsDM        bool   `json:"is_dm"`
	Status      string `json:"status"`
	System      string `json:"system"`
	PlayerCount int    `json:"player_count"`
	// The requester characters in the campaign, ordered by name; empty when the requester only runs it
	MyCharacters []CampaignCharacterRefResponse `json:"my_characters"`
}

type CampaignCharacterRefResponse struct {
	Slug string `json:"slug"`
	Name string `json:"name"`
}

type CampaignStatusResponse struct {
	Status string `json:"status"`
}

type CampaignDetailResponse struct {
	Id         int                                   `json:"id"`
	Name       string                                `json:"name"`
	Status     string                                `json:"status"`
	System     string                                `json:"system"`
	Slug       string                                `json:"slug"`
	Overview   *string                               `json:"overview"`
	IsDM       bool                                  `json:"is_dm"`
	Characters []CampaignDetailResponseCharacterItem `json:"characters"`
	InviteHash *uuid.UUID                            `json:"invite_hash"`
}

type CampaignDetailResponseCharacterItem struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	CurrentHP *int   `json:"current_hp"`
}

type CampaignListPageResponse struct {
	Items []CampaignListResponse `json:"items"`
	Total int                    `json:"total"`
}

type CampaignStatusCountsResponse struct {
	Draft    int `json:"DRAFT"`
	Active   int `json:"ACTIVE"`
	Paused   int `json:"PAUSED"`
	Archived int `json:"ARCHIVED"`
}
