package invite

type AcceptInviteRequest struct {
	CharacterSlug string `json:"character_slug"`
	// How the player is called in the campaign, 1 to 50 characters
	PlayerName string `json:"player_name"`
}

type CreateInviteRequest struct {
	CampaignID int `json:"campaign_id"`
}
