package campaign

import (
	"regexp"
	"testing"

	campaignDomain "questmaster-core/internal/campaign/domain"
	characterDomain "questmaster-core/internal/character/domain"
	characterInfra "questmaster-core/internal/character/infra/pg"
	inviteInfra "questmaster-core/internal/invite/infra/pg"
	rpgDomain "questmaster-core/internal/rpg/domain"
	"questmaster-core/internal/shared/testdb"
	userDomain "questmaster-core/internal/user/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

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

		list, err := campaigns.GetByDmId(dm)
		if err != nil || len(list) != 1 || list[0].Id != c.Id {
			t.Fatalf("expected the campaign in the DM list, got %v err=%v", list, err)
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

	expectCount := func(source string, list []campaignDomain.Campaign, err error) {
		t.Helper()
		if err != nil || len(list) != 1 || list[0].Id != campaign.Id {
			t.Fatalf("%s: expected only campaign %d, got %v err=%v", source, campaign.Id, list, err)
		}
		if list[0].PlayerCount.Value() != 2 {
			t.Fatalf("%s: expected player count 2, got %d", source, list[0].PlayerCount.Value())
		}
	}

	list, err := campaigns.GetByDmId(dm)
	expectCount("DM list", list, err)

	list, err = campaigns.GetByPlayerId(twoCharacters)
	expectCount("list of player with two characters", list, err)

	list, err = campaigns.GetByPlayerId(oneCharacter)
	expectCount("list of player with one character", list, err)

	found, err := campaigns.FindById(campaign.Id)
	if err != nil || found == nil {
		t.Fatalf("find by id: %v", err)
	}
	expectCount("find by id", []campaignDomain.Campaign{*found}, nil)

	t.Run("campaign without characters", func(t *testing.T) {
		empty := createCampaign(t, campaigns, "Empty campaign", newUser())
		found, err := campaigns.FindById(empty.Id)
		if err != nil || found == nil || found.PlayerCount.Value() != 0 {
			t.Fatalf("expected player count 0, got %v err=%v", found, err)
		}
	})
}
