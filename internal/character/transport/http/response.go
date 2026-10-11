package character

import "time"

type CharacterListResponse struct {
	Slug      string `json:"slug"`
	Name      string `json:"name"`
	System    string `json:"system"`
	CurrentHP *int   `json:"current_hp"`
	MaxHP     *int   `json:"max_hp"`
	// The requester last open of the character details, null if never opened
	LastAccessedAt *time.Time `json:"last_accessed_at"`
}

type CharacterDetailResponse struct {
	Id        int    `json:"id"`
	Slug      string `json:"slug"`
	System    string `json:"system"`
	Name      string `json:"name"`
	MaxHP     *int   `json:"max_hp"`
	CurrentHP *int   `json:"current_hp"`
	IsPlayer  bool   `json:"is_player"`
}

type CharacterCurrentHpResponse struct {
	CurrentHP int `json:"current_hp"`
}

type CharacterListPageResponse struct {
	Items []CharacterListResponse `json:"items"`
	Total int                     `json:"total"`
}
