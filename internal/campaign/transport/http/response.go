package campaign

import (
	"time"

	"github.com/google/uuid"
)

type CampaignListResponse struct {
	Slug        string `json:"slug"`
	Name        string `json:"name"`
	IsDM        bool   `json:"is_dm"`
	Status      string `json:"status"`
	System      string `json:"system"`
	PlayerCount int    `json:"player_count"`
	// The requester characters in the campaign, ordered by name; empty when the requester only runs it
	MyCharacters []CampaignCharacterRefResponse `json:"my_characters"`
	// The requester last open of the campaign details, null if never opened
	LastAccessedAt *time.Time `json:"last_accessed_at"`
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
	// Unclaimed premade characters, ordered by name; empty unless the requester is the DM
	PremadeCharacters []CampaignPremadeCharacterResponse `json:"premade_characters"`
	InviteHash        *uuid.UUID                         `json:"invite_hash"`
}

type CampaignPremadeCharacterResponse struct {
	Id        int    `json:"id"`
	Slug      string `json:"slug"`
	Name      string `json:"name"`
	CurrentHP *int   `json:"current_hp"`
	MaxHP     *int   `json:"max_hp"`
}

type CampaignDetailResponseCharacterItem struct {
	Id        int    `json:"id"`
	Slug      string `json:"slug"`
	Name      string `json:"name"`
	CurrentHP *int   `json:"current_hp"`
	MaxHP     *int   `json:"max_hp"`
	// True when the requester is the character player
	IsMine bool `json:"is_mine"`
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
