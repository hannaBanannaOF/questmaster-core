package campaign

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

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

func TestMapListReadModelToResponseLastAccessedAt(t *testing.T) {
	user := userDomain.NewUserID(uuid.New())
	opened := time.Date(2026, 10, 10, 12, 0, 0, 0, time.UTC)

	never, _ := json.Marshal(MapListReadModelToResponse(campaignDomain.Campaign{}, user))
	if !strings.Contains(string(never), `"last_accessed_at":null`) {
		t.Fatalf("expected null last_accessed_at, got %s", never)
	}

	got, _ := json.Marshal(MapListReadModelToResponse(campaignDomain.Campaign{LastAccessedAt: &opened}, user))
	if !strings.Contains(string(got), `"last_accessed_at":"2026-10-10T12:00:00Z"`) {
		t.Fatalf("expected the access time, got %s", got)
	}
}

func TestMapDetailReadModelToResponsePartyCard(t *testing.T) {
	current, max := 3, 10
	body, err := json.Marshal(MapDetailReadModelToResponse(campaignApp.CampaignDetailsReadModel{
		Characters: []campaignApp.CampaignCharacterReadModel{
			{Id: 10, Slug: "hero", Name: "Hero", CurrentHP: &current, MaxHP: &max, IsMine: true},
		},
	}))
	want := `{"id":10,"slug":"hero","name":"Hero","current_hp":3,"max_hp":10,"is_mine":true}`
	if err != nil || !strings.Contains(string(body), want) {
		t.Fatalf("expected %s in %s err=%v", want, body, err)
	}
}

func TestMapDetailReadModelToResponsePremadeCharacters(t *testing.T) {
	body, _ := json.Marshal(MapDetailReadModelToResponse(campaignApp.CampaignDetailsReadModel{}))
	if !strings.Contains(string(body), `"premade_characters":[]`) {
		t.Fatalf("expected an empty premade_characters array, got %s", body)
	}

	body, _ = json.Marshal(MapDetailReadModelToResponse(campaignApp.CampaignDetailsReadModel{
		PremadeCharacters: []campaignApp.CampaignCharacterReadModel{{Id: 20, Slug: "harvey", Name: "Harvey"}},
	}))
	if !strings.Contains(string(body), `"premade_characters":[{"id":20,"slug":"harvey","name":"Harvey","current_hp":null,"max_hp":null}]`) {
		t.Fatalf("unexpected premade_characters in %s", body)
	}
}
