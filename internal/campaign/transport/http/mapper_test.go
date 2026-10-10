package campaign

import (
	"encoding/json"
	"testing"

	campaignDomain "questmaster-core/internal/campaign/domain"
	"questmaster-core/internal/shared/pagination"
	userDomain "questmaster-core/internal/user/domain"

	"github.com/google/uuid"
)

func TestMapListPageToResponse(t *testing.T) {
	user := userDomain.NewUserID(uuid.New())

	t.Run("items and total", func(t *testing.T) {
		resp := MapListPageToResponse(pagination.Result[campaignDomain.Campaign]{
			Items: []campaignDomain.Campaign{{Name: "Lost Mine", Slug: "lost-mine", Dm: user}},
			Total: 12,
		}, user)
		if resp.Total != 12 || len(resp.Items) != 1 || resp.Items[0].Slug != "lost-mine" || !resp.Items[0].IsDM {
			t.Fatalf("unexpected response %+v", resp)
		}
	})

	t.Run("empty page renders items as an empty array", func(t *testing.T) {
		body, err := json.Marshal(MapListPageToResponse(pagination.Result[campaignDomain.Campaign]{}, user))
		if err != nil || string(body) != `{"items":[],"total":0}` {
			t.Fatalf("unexpected body %s err=%v", body, err)
		}
	})
}
