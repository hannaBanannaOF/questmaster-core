package campaign

import (
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
