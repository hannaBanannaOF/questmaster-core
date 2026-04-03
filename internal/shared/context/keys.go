package context

type contextKey string

const (
	userKey        contextKey = "User"
	campaignIDKey  contextKey = "CampaignID"
	slugKey        contextKey = "Slug"
	characterIDKey contextKey = "CharacterID"
	inviteHashKey  contextKey = "InviteHash"
)
