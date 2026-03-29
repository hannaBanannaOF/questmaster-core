package invite

type InviteDetailCharacterItem struct {
	ID   int
	Name string
}

type InviteDetailReadModel struct {
	CampaignID   int
	CampaignName string
	Characters   []InviteDetailCharacterItem
}
