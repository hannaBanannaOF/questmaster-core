package invite

import (
	"encoding/json"
	"strings"
	"testing"

	inviteApp "questmaster-core/internal/invite/app"
)

func TestMapInviteDetailsReadModelToResponseIncludesIsDM(t *testing.T) {
	for _, isDM := range []bool{true, false} {
		body, err := json.Marshal(MapInviteDetailsReadModelToResponse(inviteApp.InviteDetailReadModel{IsDM: isDM}))
		if err != nil {
			t.Fatalf("marshal: %v", err)
		}
		want := `"is_dm":false`
		if isDM {
			want = `"is_dm":true`
		}
		if !strings.Contains(string(body), want) {
			t.Fatalf("expected %s in %s", want, body)
		}
	}
}

func TestMapInviteDetailsReadModelToResponsePremadeCharacters(t *testing.T) {
	body, _ := json.Marshal(MapInviteDetailsReadModelToResponse(inviteApp.InviteDetailReadModel{}))
	if !strings.Contains(string(body), `"premade_characters":[]`) {
		t.Fatalf("expected an empty premade_characters array, got %s", body)
	}
}
