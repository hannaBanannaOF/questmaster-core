package pagination

import (
	"errors"
	"strconv"
	"strings"
	"unicode/utf8"
)

const MaxSearchLength = 100

var ErrInvalidSearch = errors.New("q must have at most 100 characters")

// NewSearch parses the raw q query value: trimmed, nil when blank, at most 100 characters.
func NewSearch(raw string) (*string, error) {
	q := strings.TrimSpace(raw)
	if q == "" {
		return nil, nil
	}
	if utf8.RuneCountInString(q) > MaxSearchLength {
		return nil, ErrInvalidSearch
	}
	return &q, nil
}

// NameContains is the SQL condition "column contains the $argIndex search", ignoring case and accents
// like the name ordering. strpos is a literal match, so % and _ in the search are not wildcards.
func NameContains(column string, argIndex int) string {
	return "strpos(unaccent(lower(" + column + ")), unaccent(lower($" + strconv.Itoa(argIndex) + "))) > 0"
}
