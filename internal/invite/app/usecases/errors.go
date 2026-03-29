package invite

import "errors"

var ErrInviteAlreadyExists = errors.New("Invite already exists for campaign")
var ErrInviteNotFound = errors.New("Invite not found")
