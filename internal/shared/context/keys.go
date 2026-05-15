package context

type contextKey string

const (
	userKey        contextKey = "User"
	filtersKey     contextKey = "Filters"
	campaignIDKey  contextKey = "CampaignID"
	slugKey        contextKey = "Slug"
	characterIDKey contextKey = "CharacterID"
	inviteHashKey  contextKey = "InviteHash"
)
