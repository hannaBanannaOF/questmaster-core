package campaign

type CampaignStatus string

const (
	StatusDraft    CampaignStatus = "DRAFT"
	StatusActive   CampaignStatus = "ACTIVE"
	StatusPaused   CampaignStatus = "PAUSED"
	StatusArchived CampaignStatus = "ARCHIVED"
)
