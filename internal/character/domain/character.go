package character

import (
	campaignDomain "questmaster-core/internal/campaign/domain"
	rpgDomain "questmaster-core/internal/rpg/domain"
)

type Character struct {
	Id         CharacterID
	Name       CharacterName
	PlayerID   rpgDomain.UserID
	System     rpgDomain.System
	CampaignID *campaignDomain.CampaignID
	Slug       rpgDomain.Slug
	Hp         *HP
}

func (c Character) IsPlayer(userID rpgDomain.UserID) bool {
	return c.PlayerID.Value() == userID.Value()
}

func (c *Character) UpdateHP(newHP HP, userID rpgDomain.UserID, campaign CampaignAccess) error {
	if err := c.CanUpdate(userID, campaign); err != nil {
		return err
	}
	c.Hp = &newHP
	return nil
}
