package character

import (
	campaignDomain "questmaster-core/internal/campaign/domain"
	userDomain "questmaster-core/internal/user/domain"
)

type CampaignAccess interface {
	IsDM(userID userDomain.UserID) bool
	IsArchived() bool
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

// CanLeaveCampaign allows the player or the campaign DM to remove the character from its campaign,
// except an unclaimed premade character or a character of an archived campaign.
func (c *Character) CanLeaveCampaign(user userDomain.UserID, campaign CampaignAccess) error {
	if !c.IsPlayer(user) && !campaign.IsDM(user) {
		return ErrNotAllowed
	}

	if c.IsPremade() {
		return ErrPremadeCannotLeave
	}

	if campaign.IsArchived() {
		return campaignDomain.ErrCampaignArchived
	}

	return nil
}

// CanDelete allows the player, and the campaign DM while the character is an unclaimed premade one.
// campaign is nil when the character has no campaign.
func (c *Character) CanDelete(user userDomain.UserID, campaign CampaignAccess) error {
	if c.IsPlayer(user) {
		return nil
	}

	if c.IsPremade() && campaign != nil && campaign.IsDM(user) {
		return nil
	}

	return ErrNotPlayer
}
