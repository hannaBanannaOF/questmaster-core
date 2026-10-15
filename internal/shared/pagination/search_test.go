package pagination

import (
	"errors"
	"strings"
	"testing"
)

func TestNewSearch(t *testing.T) {
	for _, blank := range []string{"", "   "} {
		if q, err := NewSearch(blank); err != nil || q != nil {
			t.Fatalf("NewSearch(%q) = %v, %v; want no filter", blank, q, err)
		}
	}
	if q, err := NewSearch("  Arkham  "); err != nil || q == nil || *q != "Arkham" {
		t.Fatalf("expected the trimmed search, got %v err=%v", q, err)
	}
	// 100 accented characters are 200 bytes but still valid
	if _, err := NewSearch(strings.Repeat("é", 100)); err != nil {
		t.Fatalf("expected 100 characters to be valid, got %v", err)
	}
	if _, err := NewSearch(strings.Repeat("a", 101)); !errors.Is(err, ErrInvalidSearch) {
		t.Fatalf("expected ErrInvalidSearch for 101 characters, got %v", err)
	}
}
