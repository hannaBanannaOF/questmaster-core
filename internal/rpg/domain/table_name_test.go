package rpg

import (
	"errors"
	"strings"
	"testing"
)

func TestNewTableName(t *testing.T) {
	if n, err := NewTableName("  Mestre Ana  "); err != nil || n != "Mestre Ana" {
		t.Fatalf("expected the trimmed name, got %q err=%v", n, err)
	}
	if _, err := NewTableName(strings.Repeat("é", 50)); err != nil {
		t.Fatalf("expected 50 accented characters to be valid, got %v", err)
	}
	for _, invalid := range []string{"", "   ", strings.Repeat("a", 51)} {
		if _, err := NewTableName(invalid); !errors.Is(err, ErrInvalidTableName) {
			t.Fatalf("NewTableName(%q) error = %v; want ErrInvalidTableName", invalid, err)
		}
	}
}
