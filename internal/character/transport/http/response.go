package character

import "time"

type CharacterListResponse struct {
	Slug   string `json:"slug"`
	Name   string `json:"name"`
	System string `json:"system"`
	// The requester last open of the character details, null if never opened
	LastAccessedAt *time.Time `json:"last_accessed_at"`
	// The campaign the character plays in, null when none
	Campaign *CharacterCampaignResponse `json:"campaign"`
}

type CharacterCampaignResponse struct {
	Slug string `json:"slug"`
	Name string `json:"name"`
}

type CharacterDetailResponse struct {
	Id       int    `json:"id"`
	Slug     string `json:"slug"`
	System   string `json:"system"`
	Name     string `json:"name"`
	IsPlayer bool   `json:"is_player"`
}

type CharacterListPageResponse struct {
	Items []CharacterListResponse `json:"items"`
	Total int                     `json:"total"`
}

type CharacterEditResponse struct {
	Slug       string  `json:"slug"`
	Name       string  `json:"name"`
	PlayerName *string `json:"player_name"`
}
