package character

import (
	"questmaster-core/internal/domain/campaign"
	"questmaster-core/internal/domain/rpg"
)

type Character struct {
	Id         CharacterID
	Name       CharacterName
	PlayerId   rpg.UserID
	System     rpg.System
	CampaingId *campaign.CampaignID
	Slug       rpg.Slug
	Hp         *HP
}
