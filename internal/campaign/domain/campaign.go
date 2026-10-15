package campaign

import (
	"time"

	rpgDomain "questmaster-core/internal/rpg/domain"
	userDomain "questmaster-core/internal/user/domain"
)

type Campaign struct {
	Id          CampaignID
	Name        CampaignName
	Dm          userDomain.UserID
	DmName      rpgDomain.TableName
	Status      CampaignStatus
	System      rpgDomain.System
	Slug        rpgDomain.Slug
	Overview    *CampaignOverview
	PlayerCount PlayerCount
	// The requester last open of the campaign; only set by the campaign list
	LastAccessedAt *time.Time
}

func (c Campaign) IsDM(userID userDomain.UserID) bool {
	return c.Dm.Value() == userID.Value()
}

func (c Campaign) IsArchived() bool {
	return c.Status == StatusArchived
}

func (c *Campaign) ChangeStatus(to CampaignStatus, userID userDomain.UserID) error {
	if err := c.CanEdit(userID); err != nil {
		return err
	}

	if c.Status == to {
		return nil // idempotente
	}

	if !c.Status.CanTransition(to) {
		return ErrInvalidStatusTransition
	}

	c.Status = to
	return nil
}
