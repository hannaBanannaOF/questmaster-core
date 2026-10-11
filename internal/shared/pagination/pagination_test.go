package pagination

import (
	"errors"
	"testing"
)

func TestNewPage(t *testing.T) {
	valid := []struct {
		limit, offset string
		want          Page
	}{
		{"", "", Page{Limit: DefaultLimit, Offset: 0}},
		{"1", "0", Page{Limit: 1, Offset: 0}},
		{"100", "250", Page{Limit: 100, Offset: 250}},
	}
	for _, tc := range valid {
		got, err := NewPage(tc.limit, tc.offset)
		if err != nil || got != tc.want {
			t.Fatalf("NewPage(%q, %q) = %+v, %v; want %+v", tc.limit, tc.offset, got, err, tc.want)
		}
	}

	invalid := []struct {
		limit, offset string
		want          error
	}{
		{"0", "", ErrInvalidLimit},
		{"101", "", ErrInvalidLimit},
		{"abc", "", ErrInvalidLimit},
		{"", "-1", ErrInvalidOffset},
		{"", "x", ErrInvalidOffset},
	}
	for _, tc := range invalid {
		if _, err := NewPage(tc.limit, tc.offset); !errors.Is(err, tc.want) {
			t.Fatalf("NewPage(%q, %q) error = %v; want %v", tc.limit, tc.offset, err, tc.want)
		}
	}
}

func TestNewSort(t *testing.T) {
	for raw, want := range map[string]Sort{"": SortName, "name": SortName, "recent": SortRecent} {
		if got, err := NewSort(raw); err != nil || got != want {
			t.Fatalf("NewSort(%q) = %q, %v; want %q", raw, got, err, want)
		}
	}
	for _, raw := range []string{"popular", "RECENT"} {
		if _, err := NewSort(raw); !errors.Is(err, ErrInvalidSort) {
			t.Fatalf("NewSort(%q) error = %v; want ErrInvalidSort", raw, err)
		}
	}
}
