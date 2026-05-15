package invite

type AcceptInviteRequest struct {
	CharacterSlug string `json:"character_slug"`
}

type CreateInviteRequest struct {
	CampaignID int `json:"campaign_id"`
}
