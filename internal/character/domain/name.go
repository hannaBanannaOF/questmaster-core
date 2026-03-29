package character

import (
	"strings"
)

type CharacterName string

func NewCharacterName(value string) (CharacterName, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", ErrInvalidCharacterName
	}
	return CharacterName(value), nil
}

func (n CharacterName) Value() string {
	return string(n)
}
