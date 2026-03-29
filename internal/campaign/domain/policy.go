package campaign

import rpgDomain "questmaster-core/internal/rpg/domain"

func (c *Campaign) CanEdit(userID rpgDomain.UserID) error {
	if !c.IsDM(userID) {
		return ErrNotDM
	}

	return nil
}

func (c *Campaign) CanDelete(userID rpgDomain.UserID) error {
	if !c.IsDM(userID) {
		return ErrNotDM
	}

	if c.Status != StatusDraft && c.Status != StatusArchived {
		return ErrNotDeletableStatus
	}

	return nil
}
