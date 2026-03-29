package campaign

import (
	rpg "questmaster-core/internal/rpg/domain"
	rpgDomain "questmaster-core/internal/rpg/domain"
)

type Campaign struct {
	Id       CampaignID
	Name     CampaignName
	Dm       rpgDomain.UserID
	Status   CampaignStatus
	System   rpgDomain.System
	Slug     rpgDomain.Slug
	Overview *CampaignOverview
}

func (c Campaign) IsDM(userID rpg.UserID) bool {
	return c.Dm.Value() == userID.Value()
}

func (c *Campaign) ChangeStatus(to CampaignStatus, userID rpgDomain.UserID) error {
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
