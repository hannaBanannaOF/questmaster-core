package campaign

type CreateCampaignRequest struct {
	Name     string  `json:"name"`
	Overview *string `json:"overview"`
	System   string  `json:"system"`
}

type UpdateCampaignStatusRequest struct {
	Status string `json:"status"`
}
