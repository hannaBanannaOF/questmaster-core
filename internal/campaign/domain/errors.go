package campaign

import "errors"

var ErrInvalidStatusTransition = errors.New("Invalid campaign status transition")
var ErrNotDM = errors.New("You're not this campaign DM")
var ErrNotCampaignMember = errors.New("You're not a member of this campaign")
var ErrNotDeletableStatus = errors.New("Invalid status for deletion")
var ErrEmptyCampaignName = errors.New("Empty campaign name")
var ErrInvalidCampaignStatus = errors.New("Invalid campaign status")
var ErrCampaignArchived = errors.New("Campaign is archived and doesn't take new players")
var ErrDMCannotJoin = errors.New("The DM can't join their own campaign as a player")
