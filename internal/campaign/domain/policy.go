package campaign

import userDomain "questmaster-core/internal/user/domain"

func (c *Campaign) CanEdit(userID userDomain.UserID) error {
	if !c.IsDM(userID) {
		return ErrNotDM
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
