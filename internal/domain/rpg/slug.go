package rpg

import (
	"errors"
	"regexp"
)

type Slug string

var slugRegex = regexp.MustCompile(`^[a-z0-9-]+$`)

func NewSlug(value string) (Slug, error) {
	if !slugRegex.MatchString(value) {
		return "", errors.New("Invalid slug")
	}
	return Slug(value), nil
}

func (s Slug) String() string {
	return string(s)
}
