package character

import (
	"errors"
	"strings"
)

type CharacterName string

func NewCharacterName(value string) (CharacterName, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", errors.New("Empty character name")
	}
	return CharacterName(value), nil
}

func (n CharacterName) String() string {
	return string(n)
}
