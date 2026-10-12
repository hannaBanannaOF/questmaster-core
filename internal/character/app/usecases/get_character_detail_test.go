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

type fakeDetailRepo struct {
	characterApp.CharacterRepository
	character *characterDomain.Character
	members   map[userDomain.UserID]bool
	recorded  []userDomain.UserID
	err       error
}

func (f *fakeDetailRepo) FindByID(characterDomain.CharacterID) (*characterDomain.Character, error) {
	return f.character, nil
}

func (f *fakeDetailRepo) GetByPlayerInCampaigns(user userDomain.UserID, _ []campaignDomain.CampaignID) ([]characterDomain.Character, error) {
	if f.members[user] {
		return []characterDomain.Character{{Id: 99, PlayerID: user}}, nil
	}
	return nil, nil
}

func (f *fakeDetailRepo) RecordAccess(user userDomain.UserID, _ characterDomain.CharacterID) error {
	f.recorded = append(f.recorded, user)
	return f.err
}

func TestGetCharacterDetail(t *testing.T) {
	player, dm, member, outsider := userDomain.NewUserID(uuid.New()), userDomain.NewUserID(uuid.New()),
		userDomain.NewUserID(uuid.New()), userDomain.NewUserID(uuid.New())
	campaignID := campaignDomain.NewCampaignID(1)
	finder := fakeCampaignFinder{campaign: campaignDomain.Campaign{Id: campaignID, Dm: dm}}

	inCampaign := func() *fakeDetailRepo {
		return &fakeDetailRepo{
			character: &characterDomain.Character{Id: 1, Name: "Hero", PlayerID: player, CampaignID: &campaignID},
			members:   map[userDomain.UserID]bool{member: true},
		}
	}
	read := func(repo *fakeDetailRepo, user userDomain.UserID) error {
		_, err := NewGetCharacterDetail(repo, finder).Execute(characterApp.GetCharacterDetailsCommand{ID: 1, UserID: user})
		return err
	}

	t.Run("player reads and the open is recorded", func(t *testing.T) {
		repo := inCampaign()
		if err := read(repo, player); err != nil || len(repo.recorded) != 1 || repo.recorded[0] != player {
			t.Fatalf("expected 200 with the access recorded, got err=%v recorded=%v", err, repo.recorded)
		}
	})

	for name, user := range map[string]userDomain.UserID{"campaign DM": dm, "party member": member} {
		t.Run(name+" reads without recording", func(t *testing.T) {
			repo := inCampaign()
			if err := read(repo, user); err != nil || len(repo.recorded) != 0 {
				t.Fatalf("expected 200 without recording, got err=%v recorded=%v", err, repo.recorded)
			}
		})
	}

	t.Run("non-member is forbidden and nothing is recorded", func(t *testing.T) {
		repo := inCampaign()
		if err := read(repo, outsider); !errors.Is(err, characterDomain.ErrCannotViewCharacter) || len(repo.recorded) != 0 {
			t.Fatalf("expected ErrCannotViewCharacter without recording, got err=%v recorded=%v", err, repo.recorded)
		}
	})

	t.Run("character without a campaign is only for its player", func(t *testing.T) {
		repo := &fakeDetailRepo{character: &characterDomain.Character{Id: 1, PlayerID: player}}
		if err := read(repo, outsider); !errors.Is(err, characterDomain.ErrCannotViewCharacter) {
			t.Fatalf("expected ErrCannotViewCharacter, got %v", err)
		}
	})

	t.Run("unknown character", func(t *testing.T) {
		if err := read(&fakeDetailRepo{}, player); !errors.Is(err, ErrCharacterNotFound) {
			t.Fatalf("expected ErrCharacterNotFound, got %v", err)
		}
	})

	t.Run("failing recorder still returns the character", func(t *testing.T) {
		repo := inCampaign()
		repo.err = errors.New("database is down")
		if err := read(repo, player); err != nil {
			t.Fatalf("expected the character despite the recorder error, got %v", err)
		}
	})
}
