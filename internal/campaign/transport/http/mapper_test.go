package campaign

import (
	"encoding/json"
	"strings"
	"testing"

	campaignApp "questmaster-core/internal/campaign/app"
	campaignDomain "questmaster-core/internal/campaign/domain"
	"questmaster-core/internal/shared/pagination"
	userDomain "questmaster-core/internal/user/domain"

	"github.com/google/uuid"
)

func TestMapListPageToResponse(t *testing.T) {
	user := userDomain.NewUserID(uuid.New())

	t.Run("items, total and my characters", func(t *testing.T) {
		resp := MapListPageToResponse(pagination.Result[campaignApp.CampaignListItemReadModel]{
			Items: []campaignApp.CampaignListItemReadModel{{
				Campaign:     campaignDomain.Campaign{Name: "Lost Mine", Slug: "lost-mine"},
				MyCharacters: []campaignApp.CampaignCharacterRefReadModel{{Slug: "harvey-walters", Name: "Harvey Walters"}},
			}},
			Total: 12,
		}, user)
		if resp.Total != 12 || len(resp.Items) != 1 || resp.Items[0].Slug != "lost-mine" {
			t.Fatalf("unexpected response %+v", resp)
		}
		if mine := resp.Items[0].MyCharacters; len(mine) != 1 || mine[0].Slug != "harvey-walters" || mine[0].Name != "Harvey Walters" {
			t.Fatalf("unexpected my_characters %+v", mine)
		}
	})

	t.Run("campaign without my characters renders an empty array", func(t *testing.T) {
		body, err := json.Marshal(MapListPageToResponse(pagination.Result[campaignApp.CampaignListItemReadModel]{
			Items: []campaignApp.CampaignListItemReadModel{{Campaign: campaignDomain.Campaign{Dm: user}}},
			Total: 1,
		}, user))
		if err != nil || !strings.Contains(string(body), `"my_characters":[]`) {
			t.Fatalf("expected my_characters as [], got %s err=%v", body, err)
		}
	})

	t.Run("empty page renders items as an empty array", func(t *testing.T) {
		body, err := json.Marshal(MapListPageToResponse(pagination.Result[campaignApp.CampaignListItemReadModel]{}, user))
		if err != nil || string(body) != `{"items":[],"total":0}` {
			t.Fatalf("unexpected body %s err=%v", body, err)
		}
	})
}
