package campaign

type CreateCampaignRequest struct {
	Name     string  `json:"name"`
	Overview *string `json:"overview"`
	System   string  `json:"system"`
	// How the DM is called in this campaign, 1 to 50 characters
	DmName string `json:"dm_name"`
}

type UpdateCampaignDmNameRequest struct {
	DmName string `json:"dm_name"`
}

type UpdateCampaignStatusRequest struct {
	Status string `json:"status"`
}
