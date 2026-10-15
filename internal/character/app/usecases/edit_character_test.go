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
	newRepo := func() *fakeEditRepo {
		return &fakeEditRepo{character: &characterDomain.Character{Id: 1, Name: "Harvey", PlayerID: &player, CampaignID: &campaignID}}
	}
	edit := func(repo *fakeEditRepo, user userDomain.UserID) (characterApp.EditCharacterReadModel, error) {
		return NewEditCharacter(repo, finder).Execute(characterApp.EditCharacterCommand{ID: 1, Name: name, UserID: user})
	}

	for label, user := range map[string]userDomain.UserID{"player": player, "campaign DM": dm} {
		t.Run(label+" renames", func(t *testing.T) {
			rm, err := edit(newRepo(), user)
			if err != nil || rm.Slug != "harvey-walters" || rm.Name != "Harvey Walters" {
				t.Fatalf("expected the renamed character, got %+v err=%v", rm, err)
			}
		})
	}

	t.Run("other user is forbidden and nothing is saved", func(t *testing.T) {
		repo := newRepo()
		_, err := edit(repo, userDomain.NewUserID(uuid.New()))
		if !errors.Is(err, characterDomain.ErrNotAllowed) || repo.saved != nil {
			t.Fatalf("expected ErrNotAllowed without saving, got %v saved=%v", err, repo.saved)
		}
	})

	t.Run("unknown character", func(t *testing.T) {
		if _, err := edit(&fakeEditRepo{}, player); !errors.Is(err, ErrCharacterNotFound) {
			t.Fatalf("expected ErrCharacterNotFound, got %v", err)
		}
	})
}
