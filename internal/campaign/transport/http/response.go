package campaign

import "github.com/google/uuid"

type CampaignListResponse struct {
	Slug        string `json:"slug"`
	Name        string `json:"name"`
	IsDM        bool   `json:"is_dm"`
	Status      string `json:"status"`
	System      string `json:"system"`
	PlayerCount int    `json:"player_count"`
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
	IsDM       bool                                  `json:"id_dm"`
	Characters []CampaignDetailResponseCharacterItem `json:"characters"`
	InviteHash *uuid.UUID                            `json:"invite_hash"`
}

type CampaignDetailResponseCharacterItem struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}
