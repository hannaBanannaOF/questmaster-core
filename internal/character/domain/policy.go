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

// CanView allows the player and, when the character is in a campaign, the campaign DM and its members.
// campaign is nil when the character has no campaign; isCampaignMember is resolved by the caller.
func (c *Character) CanView(user userDomain.UserID, campaign CampaignAccess, isCampaignMember bool) error {
	if c.IsPlayer(user) {
		return nil
	}

	if campaign != nil && (campaign.IsDM(user) || isCampaignMember) {
		return nil
	}

	return ErrCannotViewCharacter
}

func (c *Character) CanDelete(user userDomain.UserID) error {
	if !c.IsPlayer(user) {
		return ErrNotPlayer
	}

	return nil
}
