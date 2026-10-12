package character

import (
	"context"
	"regexp"
	"testing"

	campaignDomain "questmaster-core/internal/campaign/domain"
	campaignInfra "questmaster-core/internal/campaign/infra/pg"
	characterDomain "questmaster-core/internal/character/domain"
	rpgDomain "questmaster-core/internal/rpg/domain"
	"questmaster-core/internal/shared/pagination"
	"questmaster-core/internal/shared/testdb"
	userDomain "questmaster-core/internal/user/domain"

	"github.com/google/uuid"
)

func createCharacter(t *testing.T, r *CharacterRepositoryPG, name string, player userDomain.UserID) characterDomain.Character {
	t.Helper()
	characterName, err := characterDomain.NewCharacterName(name)
	if err != nil {
		t.Fatalf("character name: %v", err)
	}
	c, err := r.Create(characterName, player, rpgDomain.DungeonsAndDragons, nil)
	if err != nil {
		t.Fatalf("create character: %v", err)
	}
	return c
}

func TestCharacterSlugs(t *testing.T) {
	r := NewCharacterRepositoryPG(testdb.Pool(t))

	t.Run("duplicate names get a numeric suffix", func(t *testing.T) {
		// Random suffix keeps the name unique across test runs on the same database
		name := "Twin " + uuid.NewString()[:8]
		player := userDomain.NewUserID(uuid.New())
		first := createCharacter(t, r, name, player)
		second := createCharacter(t, r, name, player)
		if second.Slug.Value() != first.Slug.Value()+"-2" {
			t.Fatalf("expected %q, got %q", first.Slug.Value()+"-2", second.Slug.Value())
		}
	})

	t.Run("fallback for names without ASCII letters or digits", func(t *testing.T) {
		player := userDomain.NewUserID(uuid.New())
		c := createCharacter(t, r, "東京", player)
		if !regexp.MustCompile(`^character(-\d+)?$`).MatchString(c.Slug.Value()) {
			t.Fatalf("unexpected slug %q", c.Slug.Value())
		}

		list := listCharacters(t, r, player, characterDomain.CharacterListFilters{}, pagination.Page{Limit: pagination.DefaultLimit})
		if len(list.Items) != 1 || list.Items[0].Id != c.Id {
			t.Fatalf("expected the character in the player list, got %v", list.Items)
		}
	})
}

func TestUpdateCampaignLinksOnlyEligibleCharacters(t *testing.T) {
	db := testdb.Pool(t)
	r := NewCharacterRepositoryPG(db)
	campaigns := campaignInfra.NewCampaignRepositoryPG(db)

	newCampaign := func() campaignDomain.Campaign {
		c, err := campaigns.Create("Linking campaign", nil, userDomain.NewUserID(uuid.New()), rpgDomain.DungeonsAndDragons)
		if err != nil {
			t.Fatalf("create campaign: %v", err)
		}
		return c
	}

	t.Run("eligible character is linked", func(t *testing.T) {
		campaign := newCampaign()
		player := userDomain.NewUserID(uuid.New())
		c := createCharacter(t, r, "Eligible", player)

		linked, err := r.UpdateCampaign(campaign.Id, c.Slug, player)
		if err != nil || linked == nil || linked.CampaignID == nil || *linked.CampaignID != campaign.Id {
			t.Fatalf("expected character linked to campaign %d, got %v err=%v", campaign.Id, linked, err)
		}
	})

	t.Run("character of another player is not linked", func(t *testing.T) {
		campaign := newCampaign()
		c := createCharacter(t, r, "Not mine", userDomain.NewUserID(uuid.New()))

		linked, err := r.UpdateCampaign(campaign.Id, c.Slug, userDomain.NewUserID(uuid.New()))
		if err != nil || linked != nil {
			t.Fatalf("expected no link, got %v err=%v", linked, err)
		}
	})

	t.Run("character already in a campaign stays there", func(t *testing.T) {
		first, second := newCampaign(), newCampaign()
		player := userDomain.NewUserID(uuid.New())
		c := createCharacter(t, r, "Taken", player)
		if linked, err := r.UpdateCampaign(first.Id, c.Slug, player); err != nil || linked == nil {
			t.Fatalf("link to first campaign: %v", err)
		}

		linked, err := r.UpdateCampaign(second.Id, c.Slug, player)
		if err != nil || linked != nil {
			t.Fatalf("expected no link, got %v err=%v", linked, err)
		}
		kept, err := r.FindByID(c.Id)
		if err != nil || kept.CampaignID == nil || *kept.CampaignID != first.Id {
			t.Fatalf("expected character to stay in campaign %d, err=%v", first.Id, err)
		}
	})

	t.Run("character of another game system is not linked", func(t *testing.T) {
		campaign := newCampaign()
		player := userDomain.NewUserID(uuid.New())
		c, err := r.Create("Investigator", player, rpgDomain.CallOfCthulhu, nil)
		if err != nil {
			t.Fatalf("create character: %v", err)
		}

		linked, err := r.UpdateCampaign(campaign.Id, c.Slug, player)
		if err != nil || linked != nil {
			t.Fatalf("expected no link, got %v err=%v", linked, err)
		}
	})
}

func listCharacters(
	t *testing.T,
	r *CharacterRepositoryPG,
	player userDomain.UserID,
	filters characterDomain.CharacterListFilters,
	page pagination.Page,
) pagination.Result[characterDomain.Character] {
	t.Helper()
	list, err := r.GetAllByPlayerIDWithFilters(player, filters, page)
	if err != nil {
		t.Fatalf("list characters: %v", err)
	}
	return list
}

func characterNames(list pagination.Result[characterDomain.Character]) []string {
	out := make([]string, 0, len(list.Items))
	for _, c := range list.Items {
		out = append(out, c.Name.Value())
	}
	return out
}

func TestCharacterListOrdersAndPaginates(t *testing.T) {
	r := NewCharacterRepositoryPG(testdb.Pool(t))
	player := userDomain.NewUserID(uuid.New())
	// Created out of order, with mixed case and accents
	for _, name := range []string{"Gimli", "éowyn", "Aragorn", "boromir", "Frodo", "Denethor", "Celeborn"} {
		createCharacter(t, r, name, player)
	}
	createCharacter(t, r, "Someone else's", userDomain.NewUserID(uuid.New()))
	none := characterDomain.CharacterListFilters{}

	t.Run("preview of the first characters", func(t *testing.T) {
		list := listCharacters(t, r, player, none, pagination.Page{Limit: 5})
		want := []string{"Aragorn", "boromir", "Celeborn", "Denethor", "éowyn"}
		got := characterNames(list)
		if len(got) != len(want) || list.Total != 7 {
			t.Fatalf("expected %v of 7, got %v of %d", want, got, list.Total)
		}
		for i := range want {
			if got[i] != want[i] {
				t.Fatalf("expected %v, got %v", want, got)
			}
		}
	})

	t.Run("next page", func(t *testing.T) {
		got := characterNames(listCharacters(t, r, player, none, pagination.Page{Limit: 5, Offset: 5}))
		if len(got) != 2 || got[0] != "Frodo" || got[1] != "Gimli" {
			t.Fatalf("expected [Frodo Gimli], got %v", got)
		}
	})
}

func TestCharacterListDefaultPageSize(t *testing.T) {
	r := NewCharacterRepositoryPG(testdb.Pool(t))
	player := userDomain.NewUserID(uuid.New())
	for i := 0; i < 60; i++ {
		createCharacter(t, r, "Clone", player)
	}

	page, err := pagination.NewPage("", "")
	if err != nil {
		t.Fatalf("default page: %v", err)
	}
	list := listCharacters(t, r, player, characterDomain.CharacterListFilters{}, page)
	if len(list.Items) != 50 || list.Total != 60 {
		t.Fatalf("expected 50 items of 60, got %d of %d", len(list.Items), list.Total)
	}
}

func TestCharacterListFilters(t *testing.T) {
	db := testdb.Pool(t)
	r := NewCharacterRepositoryPG(db)
	campaigns := campaignInfra.NewCampaignRepositoryPG(db)
	player := userDomain.NewUserID(uuid.New())

	campaign, err := campaigns.Create("Filter campaign", nil, userDomain.NewUserID(uuid.New()), rpgDomain.DungeonsAndDragons)
	if err != nil {
		t.Fatalf("create campaign: %v", err)
	}
	linked := createCharacter(t, r, "Linked", player)
	if l, err := r.UpdateCampaign(campaign.Id, linked.Slug, player); err != nil || l == nil {
		t.Fatalf("link character: %v", err)
	}
	createCharacter(t, r, "Unlinked", player)
	if _, err := r.Create("Investigator", player, rpgDomain.CallOfCthulhu, nil); err != nil {
		t.Fatalf("create character: %v", err)
	}

	withoutCampaign := true
	cthulhu := rpgDomain.CallOfCthulhu
	page := pagination.Page{Limit: pagination.DefaultLimit}

	got := characterNames(listCharacters(t, r, player, characterDomain.CharacterListFilters{WithoutCampaign: &withoutCampaign}, page))
	if len(got) != 2 || got[0] != "Investigator" || got[1] != "Unlinked" {
		t.Fatalf("without_campaign: expected [Investigator Unlinked], got %v", got)
	}

	got = characterNames(listCharacters(t, r, player, characterDomain.CharacterListFilters{GameSystem: &cthulhu}, page))
	if len(got) != 1 || got[0] != "Investigator" {
		t.Fatalf("game_system: expected [Investigator], got %v", got)
	}
}

func TestGetByPlayerInCampaigns(t *testing.T) {
	db := testdb.Pool(t)
	r := NewCharacterRepositoryPG(db)
	campaigns := campaignInfra.NewCampaignRepositoryPG(db)

	dm := userDomain.NewUserID(uuid.New())
	newCampaign := func(name string) campaignDomain.Campaign {
		c, err := campaigns.Create(campaignDomain.CampaignName(name), nil, dm, rpgDomain.DungeonsAndDragons)
		if err != nil {
			t.Fatalf("create campaign: %v", err)
		}
		return c
	}
	link := func(campaign campaignDomain.Campaign, name string, player userDomain.UserID) {
		c := createCharacter(t, r, name, player)
		if l, err := r.UpdateCampaign(campaign.Id, c.Slug, player); err != nil || l == nil {
			t.Fatalf("link %s: %v", name, err)
		}
	}

	solo, crowded := newCampaign("Solo"), newCampaign("Crowded")
	player, other := userDomain.NewUserID(uuid.New()), userDomain.NewUserID(uuid.New())
	link(solo, "Harvey Walters", player)
	link(crowded, "Zed", player)
	link(crowded, "alma", player)
	link(crowded, "Not mine", other)
	ids := []campaignDomain.CampaignID{solo.Id, crowded.Id}

	got, err := r.GetByPlayerInCampaigns(player, ids)
	if err != nil {
		t.Fatalf("get characters: %v", err)
	}
	byCampaign := map[campaignDomain.CampaignID][]string{}
	for _, c := range got {
		byCampaign[*c.CampaignID] = append(byCampaign[*c.CampaignID], c.Name.Value())
	}
	if s := byCampaign[solo.Id]; len(s) != 1 || s[0] != "Harvey Walters" {
		t.Fatalf("one character: expected [Harvey Walters], got %v", s)
	}
	if c := byCampaign[crowded.Id]; len(c) != 2 || c[0] != "alma" || c[1] != "Zed" {
		t.Fatalf("several characters: expected [alma Zed], got %v", c)
	}

	t.Run("DM does not get players' characters", func(t *testing.T) {
		got, err := r.GetByPlayerInCampaigns(dm, ids)
		if err != nil || len(got) != 0 {
			t.Fatalf("expected no characters for the DM, got %v err=%v", got, err)
		}
	})
}

func TestCharacterListSortsByRecentAccess(t *testing.T) {
	db := testdb.Pool(t)
	r := NewCharacterRepositoryPG(db)
	player := userDomain.NewUserID(uuid.New())
	byName := map[string]characterDomain.Character{}
	for _, name := range []string{"Gimli", "éowyn", "Aragorn", "boromir", "Frodo", "Denethor", "Celeborn"} {
		byName[name] = createCharacter(t, r, name, player)
	}
	// Opened in this order, so Denethor is the most recent
	for _, name := range []string{"Gimli", "Frodo", "Denethor"} {
		if err := r.RecordAccess(player, byName[name].Id); err != nil {
			t.Fatalf("record access: %v", err)
		}
	}
	recent := pagination.Page{Limit: 5, Sort: pagination.SortRecent}

	t.Run("recent preview", func(t *testing.T) {
		list := listCharacters(t, r, player, characterDomain.CharacterListFilters{}, recent)
		want := []string{"Denethor", "Frodo", "Gimli", "Aragorn", "boromir"}
		got := characterNames(list)
		if len(got) != len(want) || list.Total != 7 {
			t.Fatalf("expected %v of 7, got %v of %d", want, got, list.Total)
		}
		for i := range want {
			if got[i] != want[i] {
				t.Fatalf("expected %v, got %v", want, got)
			}
		}
		if list.Items[0].LastAccessedAt == nil || list.Items[3].LastAccessedAt != nil {
			t.Fatalf("expected last_accessed_at set for opened and null for never opened")
		}
	})

	t.Run("combined with filters", func(t *testing.T) {
		withoutCampaign := true
		got := characterNames(listCharacters(t, r, player, characterDomain.CharacterListFilters{WithoutCampaign: &withoutCampaign}, recent))
		if len(got) != 5 || got[0] != "Denethor" {
			t.Fatalf("expected Denethor first among unlinked characters, got %v", got)
		}
	})

	t.Run("access removed with the character", func(t *testing.T) {
		if _, err := r.DeleteByID(byName["Denethor"].Id); err != nil {
			t.Fatalf("delete character: %v", err)
		}
		var n int
		if err := db.QueryRow(context.Background(), "SELECT COUNT(*) FROM character_access WHERE character_id = $1", byName["Denethor"].Id.Value()).Scan(&n); err != nil || n != 0 {
			t.Fatalf("expected access rows removed with the character, got %d err=%v", n, err)
		}
	})
}

func TestPremadeCharacterClaims(t *testing.T) {
	db := testdb.Pool(t)
	r := NewCharacterRepositoryPG(db)
	campaigns := campaignInfra.NewCampaignRepositoryPG(db)
	campaign, err := campaigns.Create("Premade campaign", nil, userDomain.NewUserID(uuid.New()), rpgDomain.CallOfCthulhu)
	if err != nil {
		t.Fatalf("create campaign: %v", err)
	}
	hp, _ := characterDomain.NewHP(12, 12)

	premade, err := r.CreatePremade("Harvey", campaign.Id, rpgDomain.CallOfCthulhu, &hp)
	if err != nil {
		t.Fatalf("create premade: %v", err)
	}

	t.Run("a character without a player can be read", func(t *testing.T) {
		got, err := r.FindByID(premade.Id)
		if err != nil || got == nil || !got.IsPremade() || *got.CampaignID != campaign.Id || got.Hp.Current() != 12 {
			t.Fatalf("expected the premade character back, got %+v err=%v", got, err)
		}
	})

	first, second := userDomain.NewUserID(uuid.New()), userDomain.NewUserID(uuid.New())
	claimed, err := r.ClaimPremade(premade.Id, campaign.Id, first)
	if err != nil || claimed == nil || !claimed.IsPlayer(first) {
		t.Fatalf("expected the first claim to succeed, got %+v err=%v", claimed, err)
	}

	t.Run("second claim of the same character changes nothing", func(t *testing.T) {
		again, err := r.ClaimPremade(premade.Id, campaign.Id, second)
		if err != nil || again != nil {
			t.Fatalf("expected no row changed, got %+v err=%v", again, err)
		}
		kept, _ := r.FindByID(premade.Id)
		if !kept.IsPlayer(first) {
			t.Fatalf("expected the first player to keep the character")
		}
	})

	t.Run("player who already has a character cannot claim another", func(t *testing.T) {
		another, err := r.CreatePremade("Another", campaign.Id, rpgDomain.CallOfCthulhu, nil)
		if err != nil {
			t.Fatalf("create premade: %v", err)
		}
		got, err := r.ClaimPremade(another.Id, campaign.Id, first)
		if err != nil || got != nil {
			t.Fatalf("expected no row changed, got %+v err=%v", got, err)
		}
	})
}
