package character

import (
	"time"

	campaignDomain "questmaster-core/internal/campaign/domain"
	rpgDomain "questmaster-core/internal/rpg/domain"
	userDomain "questmaster-core/internal/user/domain"
)

type Character struct {
	Id         CharacterID
	Name       CharacterName
	PlayerID   userDomain.UserID
	System     rpgDomain.System
	CampaignID *campaignDomain.CampaignID
	Slug       rpgDomain.Slug
	Hp         *HP
	// The requester last open of the character; only set by the character list
	LastAccessedAt *time.Time
}

func (c Character) IsPlayer(userID userDomain.UserID) bool {
	return c.PlayerID.Value() == userID.Value()
}

func (c *Character) UpdateHP(newHP HP, userID userDomain.UserID, campaign CampaignAccess) error {
	if err := c.CanUpdate(userID, campaign); err != nil {
		return err
	}
	c.Hp = &newHP
	return nil
}
