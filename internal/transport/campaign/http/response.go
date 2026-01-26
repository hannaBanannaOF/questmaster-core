package campaign

type CampaignListResponse struct {
	Slug   string `json:"slug"`
	Name   string `json:"name"`
	IsDM   bool   `json:"is_dm"`
	Status string `json:"status"`
	System string `json:"system"`
}

type CampaignStatusResponse struct {
	Status string `json:"status"`
}
