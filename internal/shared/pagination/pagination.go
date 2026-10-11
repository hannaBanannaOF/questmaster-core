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
var ErrInvalidSort = errors.New("sort must be name or recent")

// Sort is the order of a list. The zero value sorts by name.
type Sort string

const (
	SortName   Sort = "name"
	SortRecent Sort = "recent"
)

// NewSort parses the raw sort query value; empty uses SortName.
func NewSort(raw string) (Sort, error) {
	switch Sort(raw) {
	case "", SortName:
		return SortName, nil
	case SortRecent:
		return SortRecent, nil
	default:
		return "", ErrInvalidSort
	}
}

// OrderBy returns the ORDER BY clause for the sort: by name ignoring case and accents, then id;
// SortRecent first orders by accessed (most recent first, never accessed last).
func (s Sort) OrderBy(accessed, name, id string) string {
	order := " ORDER BY unaccent(lower(" + name + ")), " + id
	if s == SortRecent {
		order = " ORDER BY " + accessed + " DESC NULLS LAST, unaccent(lower(" + name + ")), " + id
	}
	return order
}

type Page struct {
	Limit  int
	Offset int
	Sort   Sort
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
