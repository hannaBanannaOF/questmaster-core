package character

import (
	"encoding/json"
	"testing"

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
