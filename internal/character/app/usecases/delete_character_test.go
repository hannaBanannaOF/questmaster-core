package character

import (
	"errors"
	"testing"

	campaignDomain "questmaster-core/internal/campaign/domain"
	characterApp "questmaster-core/internal/character/app"
	characterDomain "questmaster-core/internal/character/domain"
	userDomain "questmaster-core/internal/user/domain"

	"github.com/google/uuid"
)

type fakeDeleteRepo struct {
	characterApp.CharacterRepository
	character *characterDomain.Character
	deleted   bool
}

func (f *fakeDeleteRepo) FindByID(characterDomain.CharacterID) (*characterDomain.Character, error) {
	return f.character, nil
}

func (f *fakeDeleteRepo) DeleteByID(characterDomain.CharacterID) (bool, error) {
	f.deleted = true
	return true, nil
}

func TestDeleteCharacter(t *testing.T) {
	player, dm := userDomain.NewUserID(uuid.New()), userDomain.NewUserID(uuid.New())
	campaignID := campaignDomain.NewCampaignID(1)
	finder := fakeCampaignFinder{campaign: campaignDomain.Campaign{Id: campaignID, Dm: dm}}
	remove := func(c characterDomain.Character, user userDomain.UserID) (*fakeDeleteRepo, error) {
		repo := &fakeDeleteRepo{character: &c}
		return repo, NewDeleteCharacter(repo, finder).Execute(characterApp.DeleteCharacterCommand{ID: 1, UserID: user})
	}

	t.Run("DM deletes an unclaimed premade character", func(t *testing.T) {
		repo, err := remove(characterDomain.Character{Id: 1, CampaignID: &campaignID}, dm)
		if err != nil || !repo.deleted {
			t.Fatalf("expected the premade deleted, got err=%v deleted=%v", err, repo.deleted)
		}
	})

	t.Run("DM cannot delete a claimed character", func(t *testing.T) {
		repo, err := remove(characterDomain.Character{Id: 1, PlayerID: &player, CampaignID: &campaignID}, dm)
		if !errors.Is(err, characterDomain.ErrNotPlayer) || repo.deleted {
			t.Fatalf("expected ErrNotPlayer without deleting, got err=%v deleted=%v", err, repo.deleted)
		}
	})

	t.Run("other user cannot delete a premade character", func(t *testing.T) {
		repo, err := remove(characterDomain.Character{Id: 1, CampaignID: &campaignID}, userDomain.NewUserID(uuid.New()))
		if !errors.Is(err, characterDomain.ErrNotPlayer) || repo.deleted {
			t.Fatalf("expected ErrNotPlayer without deleting, got err=%v deleted=%v", err, repo.deleted)
		}
	})

	t.Run("player deletes their character", func(t *testing.T) {
		repo, err := remove(characterDomain.Character{Id: 1, PlayerID: &player}, player)
		if err != nil || !repo.deleted {
			t.Fatalf("expected the character deleted, got err=%v deleted=%v", err, repo.deleted)
		}
	})
}
