package character

import "errors"

var ErrCharacterNotFound = errors.New("Character not found")
var ErrUnavailableCharacter = errors.New("Character unavailable for requested campaign")
