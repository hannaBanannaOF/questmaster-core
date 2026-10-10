// Package pagination holds the offset pagination shared by the list endpoints.
package pagination

import (
	"errors"
	"strconv"
)

const (
	DefaultLimit = 50
	MaxLimit     = 100
)

var ErrInvalidLimit = errors.New("limit must be between 1 and 100")
var ErrInvalidOffset = errors.New("offset must be 0 or more")

type Page struct {
	Limit  int
	Offset int
}

// NewPage parses the raw limit and offset query values; empty values use the defaults.
func NewPage(rawLimit, rawOffset string) (Page, error) {
	page := Page{Limit: DefaultLimit, Offset: 0}

	if rawLimit != "" {
		limit, err := strconv.Atoi(rawLimit)
		if err != nil || limit < 1 || limit > MaxLimit {
			return Page{}, ErrInvalidLimit
		}
		page.Limit = limit
	}

	if rawOffset != "" {
		offset, err := strconv.Atoi(rawOffset)
		if err != nil || offset < 0 {
			return Page{}, ErrInvalidOffset
		}
		page.Offset = offset
	}

	return page, nil
}

// Result is one page of items and the number of items matching the filters, ignoring pagination.
type Result[T any] struct {
	Items []T
	Total int
}
