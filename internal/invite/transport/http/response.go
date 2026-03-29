package invite

type InviteDetailsCharacterListItem struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type InviteDetailsResponse struct {
	CampaignID   int                              `json:"campaign_id"`
	CampaignName string                           `json:"campaign_name"`
	Characters   []InviteDetailsCharacterListItem `json:"characters"`
}
