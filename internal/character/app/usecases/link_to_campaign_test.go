package character

import (
	"errors"
	"testing"

	campaignDomain "questmaster-core/internal/campaign/domain"
	characterApp "questmaster-core/internal/character/app"
	characterDomain "questmaster-core/internal/character/domain"
	rpgDomain "questmaster-core/internal/rpg/domain"
	userDomain "questmaster-core/internal/user/domain"

	"github.com/google/uuid"
)

type fakeLinkRepo struct {
	characterApp.CharacterRepository
	target       *characterDomain.Character
	own          []characterDomain.Character
	claimResult  *characterDomain.Character
	claimed      bool
	updateCalled bool
}

func (f *fakeLinkRepo) FindBySlug(rpgDomain.Slug) (*characterDomain.Character, error) {
	return f.target, nil
}

func (f *fakeLinkRepo) GetByPlayerInCampaigns(userDomain.UserID, []campaignDomain.CampaignID) ([]characterDomain.Character, error) {
	return f.own, nil
}

func (f *fakeLinkRepo) ClaimPremade(characterDomain.CharacterID, campaignDomain.CampaignID, userDomain.UserID) (*characterDomain.Character, error) {
	f.claimed = true
	return f.claimResult, nil
}

func (f *fakeLinkRepo) UpdateCampaign(campaignID campaignDomain.CampaignID, _ rpgDomain.Slug, user userDomain.UserID) (*characterDomain.Character, error) {
	f.updateCalled = true
	return &characterDomain.Character{Id: 1, PlayerID: &user, CampaignID: &campaignID}, nil
}

func TestLinkToCampaign(t *testing.T) {
	player := userDomain.NewUserID(uuid.New())
	campaignID, otherCampaign := campaignDomain.NewCampaignID(1), campaignDomain.NewCampaignID(2)
	premadeIn := func(id campaignDomain.CampaignID) *characterDomain.Character {
		return &characterDomain.Character{Id: 7, CampaignID: &id}
	}
	link := func(repo *fakeLinkRepo) (characterDomain.Character, error) {
		return NewLinkCharacterToCampaign(repo).LinkToCampaign(campaignID, "hero", player)
	}

	t.Run("own character keeps the existing link", func(t *testing.T) {
		repo := &fakeLinkRepo{target: &characterDomain.Character{Id: 1, PlayerID: &player}}
		if _, err := link(repo); err != nil || !repo.updateCalled || repo.claimed {
			t.Fatalf("expected the own character update, got err=%v update=%v claim=%v", err, repo.updateCalled, repo.claimed)
		}
	})

	t.Run("claims a premade character of the campaign", func(t *testing.T) {
		repo := &fakeLinkRepo{target: premadeIn(campaignID), claimResult: &characterDomain.Character{Id: 7, PlayerID: &player}}
		got, err := link(repo)
		if err != nil || !got.IsPlayer(player) || repo.updateCalled {
			t.Fatalf("expected the requester to become the player, got %+v err=%v", got, err)
		}
	})

	t.Run("player who already has a character in the campaign", func(t *testing.T) {
		repo := &fakeLinkRepo{target: premadeIn(campaignID), own: []characterDomain.Character{{Id: 3}}}
		if _, err := link(repo); !errors.Is(err, ErrAlreadyInCampaign) || repo.claimed {
			t.Fatalf("expected ErrAlreadyInCampaign without claiming, got %v claimed=%v", err, repo.claimed)
		}
	})

	t.Run("premade already claimed by someone else", func(t *testing.T) {
		repo := &fakeLinkRepo{target: premadeIn(campaignID)}
		if _, err := link(repo); !errors.Is(err, ErrPremadeAlreadyClaimed) {
			t.Fatalf("expected ErrPremadeAlreadyClaimed, got %v", err)
		}
	})

	t.Run("premade of another campaign", func(t *testing.T) {
		repo := &fakeLinkRepo{target: premadeIn(otherCampaign)}
		if _, err := link(repo); !errors.Is(err, ErrUnavailableCharacter) || repo.claimed || repo.updateCalled {
			t.Fatalf("expected ErrUnavailableCharacter without changes, got %v", err)
		}
	})
}
