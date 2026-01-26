package campaign

import (
	"questmaster-core/internal/domain/rpg"

	"github.com/google/uuid"
)

type Campaign struct {
	Id       CampaignID
	Name     CampaignName
	Dm       rpg.UserID
	Status   CampaignStatus
	System   rpg.System
	Slug     rpg.Slug
	Overview *CampaignOverview
}

func (c Campaign) IsDM(userID uuid.UUID) bool {
	return c.Dm.UUID() == userID
}

func (c *Campaign) ChangeStatus(to CampaignStatus, userId rpg.UserID) error {
	if !c.IsDM(userId.UUID()) {
		return ErrNotDM
	}

	if c.Status == to {
		return nil // idempotente
	}

	if !canTransition(c.Status, to) {
		return ErrInvalidStatusTransition
	}

	c.Status = to
	return nil
}

func canTransition(from, to CampaignStatus) bool {
	switch from {
	case StatusDraft:
		return to == StatusActive

	case StatusActive:
		return to == StatusPaused || to == StatusArchived

	case StatusPaused:
		return to == StatusActive || to == StatusArchived

	case StatusArchived:
		return false
	}

	return false
}
