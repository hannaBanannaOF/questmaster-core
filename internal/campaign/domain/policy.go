package campaign

import userDomain "questmaster-core/internal/user/domain"

func (c *Campaign) CanEdit(userID userDomain.UserID) error {
	if !c.IsDM(userID) {
		return ErrNotDM
	}

	return nil
}

// CanView allows the DM and any player with a character in the campaign.
// hasCharacter is resolved by the caller, since the campaign domain doesn't know characters.
func (c *Campaign) CanView(userID userDomain.UserID, hasCharacter bool) error {
	if c.IsDM(userID) || hasCharacter {
		return nil
	}

	return ErrNotCampaignMember
}

func (c *Campaign) CanDelete(userID userDomain.UserID) error {
	if !c.IsDM(userID) {
		return ErrNotDM
	}

	if c.Status != StatusDraft && c.Status != StatusArchived {
		return ErrNotDeletableStatus
	}

	return nil
}
