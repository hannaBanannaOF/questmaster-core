package character

import rpgDomain "questmaster-core/internal/rpg/domain"

type CampaignAccess interface {
	IsDM(userID rpgDomain.UserID) bool
}

func (c *Character) CanUpdate(
	user rpgDomain.UserID,
	campaign CampaignAccess,
) error {

	if c.IsPlayer(user) {
		return nil
	}

	if campaign != nil && campaign.IsDM(user) {
		return nil
	}

	return ErrNotAllowed
}

func (c *Character) CanDelete(user rpgDomain.UserID) error {
	if !c.IsPlayer(user) {
		return ErrNotPlayer
	}

	return nil
}
