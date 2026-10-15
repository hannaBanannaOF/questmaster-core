package rpg

import (
	"errors"
	"strings"
	"unicode/utf8"
)

const MaxTableNameLength = 50

var ErrInvalidTableName = errors.New("Table name must have 1 to 50 characters")

// TableName is how a member is called inside one campaign: the DM's name or a player's name
type TableName string

func NewTableName(value string) (TableName, error) {
	value = strings.TrimSpace(value)
	if n := utf8.RuneCountInString(value); n == 0 || n > MaxTableNameLength {
		return "", ErrInvalidTableName
	}
	return TableName(value), nil
}

func (n TableName) Value() string {
	return string(n)
}
