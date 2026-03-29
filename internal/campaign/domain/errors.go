package campaign

import "errors"

var ErrInvalidStatusTransition = errors.New("Invalid campaign status transition")
var ErrNotDM = errors.New("You're not this campaign DM")
var ErrNotDeletableStatus = errors.New("Invalid status for deletion")
var ErrEmptyCampaignName = errors.New("Empty campaign name")
var ErrInvalidCampaignStatus = errors.New("Invalid campaign status")
