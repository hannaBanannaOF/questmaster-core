package context

type contextKey string

const (
	userKey                 contextKey = "User"
	pageKey                 contextKey = "Page"
	campaignListFiltersKey  contextKey = "CampaignListFilters"
	characterListFiltersKey contextKey = "CharacterListFilters"
	campaignIDKey           contextKey = "CampaignID"
	slugKey                 contextKey = "Slug"
	characterIDKey          contextKey = "CharacterID"
	inviteHashKey           contextKey = "InviteHash"
)
