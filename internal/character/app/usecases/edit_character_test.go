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

type fakeEditRepo struct {
	characterApp.CharacterRepository
	character *characterDomain.Character
	saved     *characterDomain.Character
}

func (f *fakeEditRepo) FindByID(characterDomain.CharacterID) (*characterDomain.Character, error) {
	return f.character, nil
}

func (f *fakeEditRepo) UpdateDetails(c characterDomain.Character) (characterDomain.Character, error) {
	f.saved = &c
	c.Slug = "harvey-walters"
	return c, nil
}

func TestEditCharacter(t *testing.T) {
	player, dm := userDomain.NewUserID(uuid.New()), userDomain.NewUserID(uuid.New())
	campaignID := campaignDomain.NewCampaignID(1)
	finder := fakeCampaignFinder{campaign: campaignDomain.Campaign{Id: campaignID, Dm: dm}}
	name := characterDomain.CharacterName("Harvey Walters")
	maxHP := 13
	newRepo := func() *fakeEditRepo {
		hp, _ := characterDomain.NewHP(10, 10)
		return &fakeEditRepo{character: &characterDomain.Character{Id: 1, Name: "Harvey", PlayerID: &player, CampaignID: &campaignID, Hp: &hp}}
	}
	edit := func(repo *fakeEditRepo, cmd characterApp.EditCharacterCommand) (characterApp.EditCharacterReadModel, error) {
		return NewEditCharacter(repo, finder).Execute(cmd)
	}

	t.Run("campaign DM renames and raises max HP", func(t *testing.T) {
		repo := newRepo()
		rm, err := edit(repo, characterApp.EditCharacterCommand{ID: 1, Name: &name, MaxHP: &maxHP, UserID: dm})
		if err != nil || rm.Slug != "harvey-walters" || rm.Name != "Harvey Walters" || *rm.CurrentHP != 10 || *rm.MaxHP != 13 {
			t.Fatalf("expected the saved character, got %+v err=%v", rm, err)
		}
	})

	t.Run("other user is forbidden and nothing is saved", func(t *testing.T) {
		repo := newRepo()
		_, err := edit(repo, characterApp.EditCharacterCommand{ID: 1, Name: &name, UserID: userDomain.NewUserID(uuid.New())})
		if !errors.Is(err, characterDomain.ErrNotAllowed) || repo.saved != nil {
			t.Fatalf("expected ErrNotAllowed without saving, got %v saved=%v", err, repo.saved)
		}
	})

	t.Run("nothing to edit", func(t *testing.T) {
		repo := newRepo()
		if _, err := edit(repo, characterApp.EditCharacterCommand{ID: 1, UserID: player}); !errors.Is(err, ErrNothingToEdit) || repo.saved != nil {
			t.Fatalf("expected ErrNothingToEdit without saving, got %v", err)
		}
	})

	t.Run("invalid max HP is not saved", func(t *testing.T) {
		repo := newRepo()
		zero := 0
		if _, err := edit(repo, characterApp.EditCharacterCommand{ID: 1, MaxHP: &zero, UserID: player}); !errors.Is(err, characterDomain.ErrInvalidMaxHP) || repo.saved != nil {
			t.Fatalf("expected ErrInvalidMaxHP without saving, got %v", err)
		}
	})

	t.Run("unknown character", func(t *testing.T) {
		if _, err := edit(&fakeEditRepo{}, characterApp.EditCharacterCommand{ID: 1, Name: &name, UserID: player}); !errors.Is(err, ErrCharacterNotFound) {
			t.Fatalf("expected ErrCharacterNotFound, got %v", err)
		}
	})
}
