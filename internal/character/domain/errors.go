package character

import "errors"

var ErrNotPlayer = errors.New("You're not the player of this character sheet")
var ErrNotAllowed = errors.New("You're not allowed to make this action! You're not this character player or this character's campaign DM")
var ErrInvalidMaxHP = errors.New("Max hp must be greater than 0")
var ErrInvalidCurrentHP = errors.New("Current hp out of bounds")
var ErrInvalidCharacterName = errors.New("Invalid character name")
