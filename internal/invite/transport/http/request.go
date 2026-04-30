package invite

type AcceptInviteRequest struct {
	CharacterSheetID int `json:"character_sheet_id"`
}

type CreateInviteRequest struct {
	CampaignID int `json:"campaign_id"`
}
