package character

import userDomain "questmaster-core/internal/user/domain"

type CampaignAccess interface {
	IsDM(userID userDomain.UserID) bool
}

func (c *Character) CanUpdate(
	user userDomain.UserID,
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

func (c *Character) CanDelete(user userDomain.UserID) error {
	if !c.IsPlayer(user) {
		return ErrNotPlayer
	}

	return nil
}
