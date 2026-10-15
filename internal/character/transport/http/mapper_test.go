package character

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	characterDomain "questmaster-core/internal/character/domain"
	"questmaster-core/internal/shared/pagination"
)

func TestMapListPageToResponse(t *testing.T) {
	t.Run("items and total", func(t *testing.T) {
		resp := MapListPageToResponse(pagination.Result[characterDomain.Character]{
			Items: []characterDomain.Character{{Name: "Harvey Walters", Slug: "harvey-walters"}},
			Total: 7,
		})
		if resp.Total != 7 || len(resp.Items) != 1 || resp.Items[0].Slug != "harvey-walters" {
			t.Fatalf("unexpected response %+v", resp)
		}
	})

	t.Run("empty page renders items as an empty array", func(t *testing.T) {
		body, err := json.Marshal(MapListPageToResponse(pagination.Result[characterDomain.Character]{}))
		if err != nil || string(body) != `{"items":[],"total":0}` {
			t.Fatalf("unexpected body %s err=%v", body, err)
		}
	})
}

func TestMapListReadModelToResponseLastAccessedAt(t *testing.T) {
	opened := time.Date(2026, 10, 10, 12, 0, 0, 0, time.UTC)

	never, _ := json.Marshal(MapListReadModelToResponse(characterDomain.Character{}))
	if !strings.Contains(string(never), `"last_accessed_at":null`) {
		t.Fatalf("expected null last_accessed_at, got %s", never)
	}

	got, _ := json.Marshal(MapListReadModelToResponse(characterDomain.Character{LastAccessedAt: &opened}))
	if !strings.Contains(string(got), `"last_accessed_at":"2026-10-10T12:00:00Z"`) {
		t.Fatalf("expected the access time, got %s", got)
	}
}

func TestMapListReadModelToResponseCampaign(t *testing.T) {
	none, _ := json.Marshal(MapListReadModelToResponse(characterDomain.Character{}))
	if !strings.Contains(string(none), `"campaign":null`) {
		t.Fatalf("expected campaign null, got %s", none)
	}
	linked, _ := json.Marshal(MapListReadModelToResponse(characterDomain.Character{
		Campaign: &characterDomain.CampaignSummary{Slug: "innsmouth", Name: "Innsmouth"},
	}))
	if !strings.Contains(string(linked), `"campaign":{"slug":"innsmouth","name":"Innsmouth"}`) {
		t.Fatalf("expected the campaign, got %s", linked)
	}
}
