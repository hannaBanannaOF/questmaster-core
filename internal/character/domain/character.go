package character

import (
	"time"

	campaignDomain "questmaster-core/internal/campaign/domain"
	rpgDomain "questmaster-core/internal/rpg/domain"
	userDomain "questmaster-core/internal/user/domain"
)

type Character struct {
	Id   CharacterID
	Name CharacterName
	// Nil for a premade character the campaign DM prepared and no player claimed yet
	PlayerID *userDomain.UserID
	// How the player is called in the character's campaign; nil outside a campaign
	PlayerName *rpgDomain.TableName
	System     rpgDomain.System
	CampaignID *campaignDomain.CampaignID
	Slug       rpgDomain.Slug
	// The requester last open of the character; only set by the character list
	LastAccessedAt *time.Time
	// The linked campaign slug and name; only set by the character list
	Campaign *CampaignSummary
}

type CampaignSummary struct {
	Slug string
	Name string
}

func (c Character) IsPlayer(userID userDomain.UserID) bool {
	return c.PlayerID != nil && *c.PlayerID == userID
}

func (c Character) IsPremade() bool {
	return c.PlayerID == nil
}

// Rename changes the character name, for whoever may update it (see CanUpdate)
func (c *Character) Rename(user userDomain.UserID, campaign CampaignAccess, name CharacterName) error {
	if err := c.CanUpdate(user, campaign); err != nil {
		return err
	}

	c.Name = name
	return nil
}
