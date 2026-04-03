package user

import "strings"

type Username string

func NewUsername(value string) (Username, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", ErrInvalidUsername
	}
	return Username(value), nil
}

func (n Username) Value() string {
	return string(n)
}
