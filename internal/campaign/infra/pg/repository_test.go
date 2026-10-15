package campaign

import (
	"context"
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
	c, err := r.Create(campaignName, nil, dm, rpgDomain.DungeonsAndDragons, "Mestre")
	if err != nil {
		t.Fatalf("create campaign: %v", err)
	}
	return c
}

// createLinkedCharacter creates a character for player and links it to the campaign
func createLinkedCharacter(t *testing.T, db *pgxpool.Pool, campaignID campaignDomain.CampaignID, player userDomain.UserID) characterDomain.Character {
	t.Helper()
	r := characterInfra.NewCharacterRepositoryPG(db)
	c, err := r.Create("Hero", player, rpgDomain.DungeonsAndDragons)
	if err != nil {
		t.Fatalf("create character: %v", err)
	}
	linked, err := r.UpdateCampaign(campaignID, c.Slug, player, "Jogador")
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
	if character.PlayerName == nil || *character.PlayerName != "Jogador" {
		t.Fatalf("expected the link to store the table name, got %v", character.PlayerName)
	}
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
	if kept.PlayerName != nil {
		t.Fatalf("expected the table name cleared, got %v", *kept.PlayerName)
	}

	invite, err := invites.FindByCampaignID(campaign.Id)
	if err != nil || invite != nil {
		t.Fatalf("expected invite to be removed, got %v err=%v", invite, err)
	}
}

func TestCampaignDmName(t *testing.T) {
	db := testdb.Pool(t)
	campaigns := NewCampaignRepositoryPG(db)
	campaign := createCampaign(t, campaigns, "Named table", newUser())
	if campaign.DmName != "Mestre" {
		t.Fatalf("expected the DM name stored on create, got %q", campaign.DmName)
	}

	if err := campaigns.UpdateDmName(campaign.Id, "Mestra Bia"); err != nil {
		t.Fatalf("update DM name: %v", err)
	}
	found, err := campaigns.FindById(campaign.Id)
	if err != nil || found.DmName != "Mestra Bia" {
		t.Fatalf("expected the new DM name, got %+v err=%v", found, err)
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

func TestCountByStatusForUser(t *testing.T) {
	db := testdb.Pool(t)
	campaigns := NewCampaignRepositoryPG(db)
	user := newUser()

	// Runs 2 ACTIVE and 1 DRAFT, plays in 1 ACTIVE; another user's campaign must not count
	for _, name := range []string{"Active one", "Active two"} {
		c := createCampaign(t, campaigns, name, user)
		if _, err := campaigns.UpdateStatus(campaignDomain.StatusActive, c.Id); err != nil {
			t.Fatalf("activate: %v", err)
		}
	}
	createCampaign(t, campaigns, "Draft", user)
	played := createCampaign(t, campaigns, "Played", newUser())
	createLinkedCharacter(t, db, played.Id, user)
	if _, err := campaigns.UpdateStatus(campaignDomain.StatusActive, played.Id); err != nil {
		t.Fatalf("activate: %v", err)
	}
	createCampaign(t, campaigns, "Someone else's", newUser())

	dm, player := campaignDomain.RoleDM, campaignDomain.RolePlayer
	cases := []struct {
		name    string
		filters campaignDomain.CampaignListFilters
		want    map[campaignDomain.CampaignStatus]int
	}{
		{"role dm", campaignDomain.CampaignListFilters{Role: &dm},
			map[campaignDomain.CampaignStatus]int{campaignDomain.StatusActive: 2, campaignDomain.StatusDraft: 1}},
		{"role player", campaignDomain.CampaignListFilters{Role: &player},
			map[campaignDomain.CampaignStatus]int{campaignDomain.StatusActive: 1}},
		{"both roles", campaignDomain.CampaignListFilters{},
			map[campaignDomain.CampaignStatus]int{campaignDomain.StatusActive: 3, campaignDomain.StatusDraft: 1}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := campaigns.CountByStatusForUser(user, tc.filters)
			if err != nil || len(got) != len(tc.want) {
				t.Fatalf("expected %v, got %v err=%v", tc.want, got, err)
			}
			for status, n := range tc.want {
				if got[status] != n {
					t.Fatalf("expected %v, got %v", tc.want, got)
				}
			}
		})
	}
}

func countAccessRows(t *testing.T, db *pgxpool.Pool, campaignID campaignDomain.CampaignID) int {
	t.Helper()
	var n int
	if err := db.QueryRow(context.Background(), "SELECT COUNT(*) FROM campaign_access WHERE campaign_id = $1", campaignID.Value()).Scan(&n); err != nil {
		t.Fatalf("count access rows: %v", err)
	}
	return n
}

func TestRecordCampaignAccess(t *testing.T) {
	db := testdb.Pool(t)
	campaigns := NewCampaignRepositoryPG(db)
	dm := newUser()
	campaign := createCampaign(t, campaigns, "Accessed", dm)
	recent := pagination.Page{Limit: pagination.DefaultLimit, Sort: pagination.SortRecent}

	if err := campaigns.RecordAccess(dm, campaign.Id); err != nil {
		t.Fatalf("first access: %v", err)
	}
	first := listFor(t, campaigns, dm, campaignDomain.CampaignListFilters{}, recent).Items[0].LastAccessedAt
	if first == nil {
		t.Fatalf("expected last_accessed_at after the first access")
	}

	if err := campaigns.RecordAccess(dm, campaign.Id); err != nil {
		t.Fatalf("repeated access: %v", err)
	}
	second := listFor(t, campaigns, dm, campaignDomain.CampaignListFilters{}, recent).Items[0].LastAccessedAt
	if n := countAccessRows(t, db, campaign.Id); n != 1 || second == nil || second.Before(*first) {
		t.Fatalf("expected one row moved forward, got %d rows, %v then %v", n, first, second)
	}

	if _, err := campaigns.DeleteById(campaign.Id); err != nil {
		t.Fatalf("delete campaign: %v", err)
	}
	if n := countAccessRows(t, db, campaign.Id); n != 0 {
		t.Fatalf("expected access rows removed with the campaign, got %d", n)
	}
}

func TestListForUserSortsByRecentAccess(t *testing.T) {
	db := testdb.Pool(t)
	campaigns := NewCampaignRepositoryPG(db)
	dm, otherMember := newUser(), newUser()
	byName := map[string]campaignDomain.Campaign{}
	for _, name := range []string{"Beta", "alfa", "Ômega"} {
		byName[name] = createCampaign(t, campaigns, name, dm)
	}
	createLinkedCharacter(t, db, byName["Ômega"].Id, otherMember)

	// "Beta" first, then "alfa"; "Ômega" only opened by another member
	for _, open := range []struct {
		user userDomain.UserID
		name string
	}{{dm, "Beta"}, {dm, "alfa"}, {otherMember, "Ômega"}} {
		if err := campaigns.RecordAccess(open.user, byName[open.name].Id); err != nil {
			t.Fatalf("record access: %v", err)
		}
	}

	list := listFor(t, campaigns, dm, campaignDomain.CampaignListFilters{}, pagination.Page{Limit: 10, Sort: pagination.SortRecent})
	if got := names(list); len(got) != 3 || got[0] != "alfa" || got[1] != "Beta" || got[2] != "Ômega" {
		t.Fatalf("expected [alfa Beta Ômega], got %v", got)
	}
	if list.Items[0].LastAccessedAt == nil || list.Items[2].LastAccessedAt != nil {
		t.Fatalf("expected last_accessed_at set for opened and null for never opened")
	}

	t.Run("name sort keeps last_accessed_at", func(t *testing.T) {
		list := listFor(t, campaigns, dm, campaignDomain.CampaignListFilters{}, pagination.Page{Limit: 10})
		if got := names(list); got[0] != "alfa" || list.Items[0].LastAccessedAt == nil {
			t.Fatalf("expected alfa first with its access, got %v", got)
		}
	})

	t.Run("combined with filters", func(t *testing.T) {
		dmRole, draft := campaignDomain.RoleDM, campaignDomain.StatusDraft
		list := listFor(t, campaigns, dm, campaignDomain.CampaignListFilters{Role: &dmRole, Status: &draft},
			pagination.Page{Limit: 1, Sort: pagination.SortRecent})
		if got := names(list); len(got) != 1 || got[0] != "alfa" || list.Total != 3 {
			t.Fatalf("expected [alfa] of 3, got %v of %d", got, list.Total)
		}
	})
}

func TestPremadeCharactersInCampaign(t *testing.T) {
	db := testdb.Pool(t)
	campaigns := NewCampaignRepositoryPG(db)
	characters := characterInfra.NewCharacterRepositoryPG(db)
	dm, player, claimer := newUser(), newUser(), newUser()
	campaign := createCampaign(t, campaigns, "Premade table", dm)

	createLinkedCharacter(t, db, campaign.Id, player)
	claimed, err := characters.CreatePremade("Claimed", campaign.Id, rpgDomain.DungeonsAndDragons)
	if err != nil {
		t.Fatalf("create premade: %v", err)
	}
	if c, err := characters.ClaimPremade(claimed.Id, campaign.Id, claimer, "Jogador"); err != nil || c == nil {
		t.Fatalf("claim premade: %v", err)
	}
	unclaimed, err := characters.CreatePremade("Unclaimed", campaign.Id, rpgDomain.DungeonsAndDragons)
	if err != nil {
		t.Fatalf("create premade: %v", err)
	}

	t.Run("party has players only, premades list has unclaimed only", func(t *testing.T) {
		party, err := characters.GetAllByCampaignID(campaign.Id)
		if err != nil || len(party) != 2 {
			t.Fatalf("expected the player and the claimer in the party, got %v err=%v", party, err)
		}
		premades, err := characters.GetPremadeByCampaignID(campaign.Id)
		if err != nil || len(premades) != 1 || premades[0].Id != unclaimed.Id {
			t.Fatalf("expected only the unclaimed premade, got %v err=%v", premades, err)
		}
	})

	t.Run("unclaimed premades do not count as players", func(t *testing.T) {
		found, err := campaigns.FindById(campaign.Id)
		if err != nil || found.PlayerCount.Value() != 2 {
			t.Fatalf("expected player count 2, got %v err=%v", found, err)
		}
	})

	t.Run("premades are not in the DM character list nor my_characters", func(t *testing.T) {
		list, err := characters.GetAllByPlayerIDWithFilters(dm, characterDomain.CharacterListFilters{}, firstPage)
		if err != nil || list.Total != 0 {
			t.Fatalf("expected no characters for the DM, got %v err=%v", list.Items, err)
		}
		mine, err := characters.GetByPlayerInCampaigns(dm, []campaignDomain.CampaignID{campaign.Id})
		if err != nil || len(mine) != 0 {
			t.Fatalf("expected no my_characters for the DM, got %v err=%v", mine, err)
		}
	})

	t.Run("deleting the campaign deletes only the unclaimed premade", func(t *testing.T) {
		if deleted, err := campaigns.DeleteById(campaign.Id); err != nil || !deleted {
			t.Fatalf("delete campaign: deleted=%v err=%v", deleted, err)
		}
		gone, err := characters.FindByID(unclaimed.Id)
		if err != nil || gone != nil {
			t.Fatalf("expected the unclaimed premade deleted, got %v err=%v", gone, err)
		}
		kept, err := characters.FindByID(claimed.Id)
		if err != nil || kept == nil || !kept.IsPlayer(claimer) || kept.CampaignID != nil {
			t.Fatalf("expected the claimed character kept and unlinked, got %+v err=%v", kept, err)
		}
	})
}

func TestRemovedCharacterNoLongerGrantsMembership(t *testing.T) {
	db := testdb.Pool(t)
	campaigns := NewCampaignRepositoryPG(db)
	characters := characterInfra.NewCharacterRepositoryPG(db)
	dm, leaving, staying := newUser(), newUser(), newUser()
	campaign := createCampaign(t, campaigns, "Leaving table", dm)
	character := createLinkedCharacter(t, db, campaign.Id, leaving)
	createLinkedCharacter(t, db, campaign.Id, staying)

	if unlinked, err := characters.UnlinkFromCampaign(character.Id, campaign.Id); err != nil || !unlinked {
		t.Fatalf("unlink: unlinked=%v err=%v", unlinked, err)
	}

	if list := listFor(t, campaigns, leaving, campaignDomain.CampaignListFilters{}, firstPage); list.Total != 0 {
		t.Fatalf("expected the campaign out of the leaving player list, got %v", names(list))
	}
	mine, err := characters.GetByPlayerInCampaigns(leaving, []campaignDomain.CampaignID{campaign.Id})
	if err != nil || len(mine) != 0 {
		t.Fatalf("expected no membership left, got %v err=%v", mine, err)
	}
	found, err := campaigns.FindById(campaign.Id)
	if err != nil || found.PlayerCount.Value() != 1 {
		t.Fatalf("expected player count 1, got %v err=%v", found, err)
	}
	kept, err := characters.FindByID(character.Id)
	if err != nil || kept == nil || !kept.IsPlayer(leaving) || kept.CampaignID != nil || kept.PlayerName != nil {
		t.Fatalf("expected the character kept by its player with no campaign, got %+v err=%v", kept, err)
	}

	t.Run("unlinking again changes nothing", func(t *testing.T) {
		if unlinked, err := characters.UnlinkFromCampaign(character.Id, campaign.Id); err != nil || unlinked {
			t.Fatalf("expected no row changed, got unlinked=%v err=%v", unlinked, err)
		}
	})
}

func TestListForUserSearchAndWithoutPlayers(t *testing.T) {
	db := testdb.Pool(t)
	campaigns := NewCampaignRepositoryPG(db)
	dm := newUser()
	// Unique suffix so names from other runs on the same database never match
	tag := uuid.NewString()[:8]
	alem := createCampaign(t, campaigns, "Além das Montanhas "+tag, dm)
	createCampaign(t, campaigns, "Arkham, 1925 "+tag, dm)
	createLinkedCharacter(t, db, alem.Id, newUser())
	createCampaign(t, campaigns, "Além do outro "+tag, newUser())
	search := func(q string) *string { return &q }

	t.Run("search ignores case and accents", func(t *testing.T) {
		list := listFor(t, campaigns, dm, campaignDomain.CampaignListFilters{Search: search("ALEM DAS")}, firstPage)
		if got := names(list); len(got) != 1 || got[0] != "Além das Montanhas "+tag || list.Total != 1 {
			t.Fatalf("expected only Além das Montanhas, got %v (total %d)", got, list.Total)
		}
	})

	t.Run("other users' campaigns are not found", func(t *testing.T) {
		list := listFor(t, campaigns, dm, campaignDomain.CampaignListFilters{Search: search("outro " + tag)}, firstPage)
		if list.Total != 0 {
			t.Fatalf("expected no campaign, got %v", names(list))
		}
	})

	t.Run("search is literal", func(t *testing.T) {
		list := listFor(t, campaigns, dm, campaignDomain.CampaignListFilters{Search: search("%")}, firstPage)
		if list.Total != 0 {
			t.Fatalf("expected %% to match nothing, got %v", names(list))
		}
	})

	t.Run("without players combines with role and search", func(t *testing.T) {
		dmRole := campaignDomain.RoleDM
		list := listFor(t, campaigns, dm, campaignDomain.CampaignListFilters{Role: &dmRole, WithoutPlayers: true, Search: search(tag)}, firstPage)
		if got := names(list); len(got) != 1 || got[0] != "Arkham, 1925 "+tag || list.Total != 1 {
			t.Fatalf("expected only the campaign nobody joined, got %v (total %d)", got, list.Total)
		}
	})
}
