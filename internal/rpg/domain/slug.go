package rpg

import (
	"regexp"
)

type Slug string

var slugRegex = regexp.MustCompile(`^[a-z0-9-]+$`)

func NewSlug(value string) (Slug, error) {
	if !slugRegex.MatchString(value) {
		return "", ErrInvalidSlug
	}
	return Slug(value), nil
}

func (s Slug) Value() string {
	return string(s)
}
