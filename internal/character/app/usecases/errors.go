package character

import "errors"

var ErrCharacterNotFound = errors.New("Character not found")
var ErrAlreadyInCampaign = errors.New("You already have a character in this campaign")
var ErrPremadeAlreadyClaimed = errors.New("This premade character was already claimed by another player")
var ErrUnavailableCharacter = errors.New("Character unavailable for requested campaign")
