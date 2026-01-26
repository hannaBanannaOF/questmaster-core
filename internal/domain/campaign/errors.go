package campaign

import "errors"

var ErrInvalidStatusTransition = errors.New("Invalid campaign status transition")
var ErrNotDM = errors.New("Unable to update campaign! You're not the DM")
