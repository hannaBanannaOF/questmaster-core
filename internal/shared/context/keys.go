package context

type contextKey string

const (
	userIDKey      contextKey = "UserID"
	campaignIDKey  contextKey = "CampaignID"
	slugKey        contextKey = "Slug"
	characterIDKey contextKey = "CharacterID"
	inviteHashKey  contextKey = "InviteHash"
)
