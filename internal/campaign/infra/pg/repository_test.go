package campaign

import (
	"regexp"
	"testing"

	campaignDomain "questmaster-core/internal/campaign/domain"
	characterDomain "questmaster-core/internal/character/domain"
	characterInfra "questmaster-core/internal/character/infra/pg"
	inviteInfra "questmaster-core/internal/invite/infra/pg"
	rpgDomain "questmaster-core/internal/rpg/domain"
	"questmaster-core/internal/shared/pagination"
	"questmaster-core/internal/shared/testdb"
	userDomain "questmaster-core/internal/user/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

var firstPage = pagination.Page{Limit: pagination.DefaultLimit}

func listFor(
	t *testing.T,
	r *CampaignRepositoryPG,
	user userDomain.UserID,
	filters campaignDomain.CampaignListFilters,
	page pagination.Page,
) pagination.Result[campaignDomain.Campaign] {
	t.Helper()
	list, err := r.ListForUser(user, filters, page)
	if err != nil {
		t.Fatalf("list campaigns: %v", err)
	}
	return list
}

func newUser() userDomain.UserID {
	return userDomain.NewUserID(uuid.New())
}

func createCampaign(t *testing.T, r *CampaignRepositoryPG, name string, dm userDomain.UserID) campaignDomain.Campaign {
	t.Helper()
	campaignName, err := campaignDomain.NewCampaignName(name)
	if err != nil {
		t.Fatalf("campaign name: %v", err)
	}
	c, err := r.Create(campaignName, nil, dm, rpgDomain.DungeonsAndDragons)
	if err != nil {
		t.Fatalf("create campaign: %v", err)
	}
	return c
}

// createLinkedCharacter creates a character for player and links it to the campaign
func createLinkedCharacter(t *testing.T, db *pgxpool.Pool, campaignID campaignDomain.CampaignID, player userDomain.UserID) characterDomain.Character {
	t.Helper()
	r := characterInfra.NewCharacterRepositoryPG(db)
	c, err := r.Create("Hero", player, rpgDomain.DungeonsAndDragons, nil)
	if err != nil {
		t.Fatalf("create character: %v", err)
	}
	linked, err := r.UpdateCampaign(campaignID, c.Slug, player)
	if err != nil || linked == nil {
		t.Fatalf("link character: %v", err)
	}
	return *linked
}

func TestCampaignSlugs(t *testing.T) {
	db := testdb.Pool(t)
	campaigns := NewCampaignRepositoryPG(db)

	t.Run("derived from name without accents and punctuation", func(t *testing.T) {
		c := createCampaign(t, campaigns, "Mina Assombrada Ébria!", newUser())
		if !regexp.MustCompile(`^mina-assombrada-ebria(-\d+)?$`).MatchString(c.Slug.Value()) {
			t.Fatalf("unexpected slug %q", c.Slug.Value())
		}
	})

	t.Run("fallback for names without ASCII letters or digits", func(t *testing.T) {
		dm := newUser()
		c := createCampaign(t, campaigns, "???", dm)
		if !regexp.MustCompile(`^campaign(-\d+)?$`).MatchString(c.Slug.Value()) {
			t.Fatalf("unexpected slug %q", c.Slug.Value())
		}

		list := listFor(t, campaigns, dm, campaignDomain.CampaignListFilters{}, firstPage)
		if len(list.Items) != 1 || list.Items[0].Id != c.Id {
			t.Fatalf("expected the campaign in the DM list, got %v", list.Items)
		}
	})
}

func TestDeleteCampaignKeepsCharactersAndRemovesInvite(t *testing.T) {
	db := testdb.Pool(t)
	campaigns := NewCampaignRepositoryPG(db)
	characters := characterInfra.NewCharacterRepositoryPG(db)
	invites := inviteInfra.NewInviteRepositoryPG(db)

	campaign := createCampaign(t, campaigns, "Doomed campaign", newUser())
	player := newUser()
	character := createLinkedCharacter(t, db, campaign.Id, player)
	if _, err := invites.Create(campaign.Id); err != nil {
		t.Fatalf("create invite: %v", err)
	}

	deleted, err := campaigns.DeleteById(campaign.Id)
	if err != nil || !deleted {
		t.Fatalf("expected campaign to be deleted, deleted=%v err=%v", deleted, err)
	}

	kept, err := characters.FindByID(character.Id)
	if err != nil || kept == nil {
		t.Fatalf("expected character to still exist, err=%v", err)
	}
	if kept.CampaignID != nil {
		t.Fatalf("expected character to be unlinked, got campaign %v", *kept.CampaignID)
	}
	if !kept.IsPlayer(player) {
		t.Fatalf("expected character to keep its player")
	}

	invite, err := invites.FindByCampaignID(campaign.Id)
	if err != nil || invite != nil {
		t.Fatalf("expected invite to be removed, got %v err=%v", invite, err)
	}
}

func TestPlayerCountCountsDistinctPlayers(t *testing.T) {
	db := testdb.Pool(t)
	campaigns := NewCampaignRepositoryPG(db)

	dm := newUser()
	campaign := createCampaign(t, campaigns, "Crowded campaign", dm)
	twoCharacters, oneCharacter := newUser(), newUser()
	createLinkedCharacter(t, db, campaign.Id, twoCharacters)
	createLinkedCharacter(t, db, campaign.Id, twoCharacters)
	createLinkedCharacter(t, db, campaign.Id, oneCharacter)

	expectCount := func(source string, list []campaignDomain.Campaign) {
		t.Helper()
		if len(list) != 1 || list[0].Id != campaign.Id {
			t.Fatalf("%s: expected only campaign %d, got %v", source, campaign.Id, list)
		}
		if list[0].PlayerCount.Value() != 2 {
			t.Fatalf("%s: expected player count 2, got %d", source, list[0].PlayerCount.Value())
		}
	}

	all := campaignDomain.CampaignListFilters{}
	expectCount("DM list", listFor(t, campaigns, dm, all, firstPage).Items)
	expectCount("list of player with two characters", listFor(t, campaigns, twoCharacters, all, firstPage).Items)
	expectCount("list of player with one character", listFor(t, campaigns, oneCharacter, all, firstPage).Items)

	found, err := campaigns.FindById(campaign.Id)
	if err != nil || found == nil {
		t.Fatalf("find by id: %v", err)
	}
	expectCount("find by id", []campaignDomain.Campaign{*found})

	t.Run("campaign without characters", func(t *testing.T) {
		empty := createCampaign(t, campaigns, "Empty campaign", newUser())
		found, err := campaigns.FindById(empty.Id)
		if err != nil || found == nil || found.PlayerCount.Value() != 0 {
			t.Fatalf("expected player count 0, got %v err=%v", found, err)
		}
	})
}

func names(list pagination.Result[campaignDomain.Campaign]) []string {
	out := make([]string, 0, len(list.Items))
	for _, c := range list.Items {
		out = append(out, c.Name.Value())
	}
	return out
}

func TestListForUserOrdersAndPaginates(t *testing.T) {
	campaigns := NewCampaignRepositoryPG(testdb.Pool(t))
	dm := newUser()
	for _, name := range []string{"Beta", "alfa", "Ômega"} {
		createCampaign(t, campaigns, name, dm)
	}
	all := campaignDomain.CampaignListFilters{}

	first := listFor(t, campaigns, dm, all, pagination.Page{Limit: 2})
	if got := names(first); len(got) != 2 || got[0] != "alfa" || got[1] != "Beta" || first.Total != 3 {
		t.Fatalf("first page: expected [alfa Beta] of 3, got %v of %d", got, first.Total)
	}

	next := listFor(t, campaigns, dm, all, pagination.Page{Limit: 2, Offset: 2})
	if got := names(next); len(got) != 1 || got[0] != "Ômega" || next.Total != 3 {
		t.Fatalf("next page: expected [Ômega] of 3, got %v of %d", got, next.Total)
	}
}

func TestListForUserFilters(t *testing.T) {
	db := testdb.Pool(t)
	campaigns := NewCampaignRepositoryPG(db)
	user := newUser()

	// user runs a DRAFT campaign and plays in an ACTIVE one; a third campaign is unrelated
	createCampaign(t, campaigns, "Mastered", user)
	played := createCampaign(t, campaigns, "Played", newUser())
	createLinkedCharacter(t, db, played.Id, user)
	if _, err := campaigns.UpdateStatus(campaignDomain.StatusActive, played.Id); err != nil {
		t.Fatalf("activate campaign: %v", err)
	}
	createCampaign(t, campaigns, "Someone else's", newUser())

	dm, player := campaignDomain.RoleDM, campaignDomain.RolePlayer
	active, draft := campaignDomain.StatusActive, campaignDomain.StatusDraft

	cases := []struct {
		name    string
		filters campaignDomain.CampaignListFilters
		want    []string
	}{
		{"both roles", campaignDomain.CampaignListFilters{}, []string{"Mastered", "Played"}},
		{"role dm", campaignDomain.CampaignListFilters{Role: &dm}, []string{"Mastered"}},
		{"role player", campaignDomain.CampaignListFilters{Role: &player}, []string{"Played"}},
		{"player and active", campaignDomain.CampaignListFilters{Role: &player, Status: &active}, []string{"Played"}},
		{"player and draft", campaignDomain.CampaignListFilters{Role: &player, Status: &draft}, []string{}},
		{"status draft", campaignDomain.CampaignListFilters{Status: &draft}, []string{"Mastered"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			list := listFor(t, campaigns, user, tc.filters, firstPage)
			got := names(list)
			if len(got) != len(tc.want) || list.Total != len(tc.want) {
				t.Fatalf("expected %v (total %d), got %v (total %d)", tc.want, len(tc.want), got, list.Total)
			}
			for i := range got {
				if got[i] != tc.want[i] {
					t.Fatalf("expected %v, got %v", tc.want, got)
				}
			}
		})
	}
}
