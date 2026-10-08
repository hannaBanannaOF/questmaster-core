package character

import (
	"regexp"
	"testing"

	campaignDomain "questmaster-core/internal/campaign/domain"
	campaignInfra "questmaster-core/internal/campaign/infra/pg"
	characterDomain "questmaster-core/internal/character/domain"
	rpgDomain "questmaster-core/internal/rpg/domain"
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

		list, err := r.GetAllByPlayerIDWithFilters(player, &characterDomain.CharacterListFilters{})
		if err != nil || len(list) != 1 || list[0].Id != c.Id {
			t.Fatalf("expected the character in the player list, got %v err=%v", list, err)
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
