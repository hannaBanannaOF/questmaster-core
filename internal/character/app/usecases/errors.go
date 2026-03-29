package character

import "errors"

var ErrCharacterNotFound = errors.New("Character not found")
var ErrAlreadyEnrolled = errors.New("Character already enrolled in a different campaign")
