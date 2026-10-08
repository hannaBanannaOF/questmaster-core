package httperrors

import (
	"errors"
	"net/http"
	"strings"
	"testing"
)

func TestFromHidesInternalErrors(t *testing.T) {
	httpErr := From(errors.New(`pq: relation "campaign" does not exist`))

	if httpErr.Status != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", httpErr.Status)
	}
	if strings.Contains(httpErr.Message, "campaign") {
		t.Fatalf("internal error leaked to client: %q", httpErr.Message)
	}
}
