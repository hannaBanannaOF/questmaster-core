package invite

import "github.com/google/uuid"

type InviteDetailReadModel struct {
	InviteHash          uuid.UUID
	CampaignSlug        string
	CampaignName        string
	CampaignSystem      string
	CampaignOverview    *string
	CampaignPlayerCount int
}
