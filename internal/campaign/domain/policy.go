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

// CanInvite allows the DM to create an invite while the campaign is not archived.
func (c *Campaign) CanInvite(userID userDomain.UserID) error {
	if err := c.CanEdit(userID); err != nil {
		return err
	}

	if c.Status == StatusArchived {
		return ErrCampaignArchived
	}

	return nil
}

// CanJoin allows anyone but the DM to join through an invite while the campaign is not archived.
func (c *Campaign) CanJoin(userID userDomain.UserID) error {
	if c.Status == StatusArchived {
		return ErrCampaignArchived
	}

	if c.IsDM(userID) {
		return ErrDMCannotJoin
	}

	return nil
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
