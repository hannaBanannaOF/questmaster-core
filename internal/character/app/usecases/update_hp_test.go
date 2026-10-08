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

type fakeCharacterRepo struct {
	characterApp.CharacterRepository
	character *characterDomain.Character
}

func (f fakeCharacterRepo) FindByID(characterDomain.CharacterID) (*characterDomain.Character, error) {
	return f.character, nil
}

func (f fakeCharacterRepo) UpdateHP(newHP characterDomain.HP, _ characterDomain.CharacterID) (characterDomain.Character, error) {
	updated := *f.character
	updated.Hp = &newHP
	return updated, nil
}

type fakeCampaignFinder struct {
	campaign campaignDomain.Campaign
}

func (f fakeCampaignFinder) FindByID(campaignDomain.CampaignID) (campaignDomain.Campaign, error) {
	return f.campaign, nil
}

func TestUpdateHP(t *testing.T) {
	player := userDomain.NewUserID(uuid.New())
	dm := userDomain.NewUserID(uuid.New())
	campaignID := campaignDomain.NewCampaignID(1)
	finder := fakeCampaignFinder{campaign: campaignDomain.Campaign{Id: campaignID, Dm: dm}}

	newCharacterRepo := func() fakeCharacterRepo {
		hp, _ := characterDomain.NewHP(10, 10)
		return fakeCharacterRepo{character: &characterDomain.Character{
			Id: 1, PlayerID: player, CampaignID: &campaignID, Hp: &hp,
		}}
	}

	allowed := map[string]userDomain.UserID{"player": player, "campaign DM": dm}
	for name, user := range allowed {
		t.Run(name+" updates HP", func(t *testing.T) {
			rm, err := NewUpdateHP(newCharacterRepo(), finder).Execute(characterApp.UpdateHPCommand{ID: 1, NewHP: 4, UserID: user})
			if err != nil || rm.CurrentHP != 4 {
				t.Fatalf("expected current hp 4, got %d err=%v", rm.CurrentHP, err)
			}
		})
	}

	t.Run("other user is forbidden", func(t *testing.T) {
		other := userDomain.NewUserID(uuid.New())
		_, err := NewUpdateHP(newCharacterRepo(), finder).Execute(characterApp.UpdateHPCommand{ID: 1, NewHP: 4, UserID: other})
		if !errors.Is(err, characterDomain.ErrNotAllowed) {
			t.Fatalf("expected ErrNotAllowed, got %v", err)
		}
	})

	for _, value := range []int{11, -1} {
		t.Run("out of bounds hp is rejected", func(t *testing.T) {
			_, err := NewUpdateHP(newCharacterRepo(), finder).Execute(characterApp.UpdateHPCommand{ID: 1, NewHP: value, UserID: player})
			if !errors.Is(err, characterDomain.ErrInvalidCurrentHP) {
				t.Fatalf("expected ErrInvalidCurrentHP for %d, got %v", value, err)
			}
		})
	}

	t.Run("character without hp is rejected", func(t *testing.T) {
		repo := fakeCharacterRepo{character: &characterDomain.Character{Id: 1, PlayerID: player}}
		_, err := NewUpdateHP(repo, nil).Execute(characterApp.UpdateHPCommand{ID: 1, NewHP: 5, UserID: player})
		if !errors.Is(err, characterDomain.ErrCharacterWithoutHP) {
			t.Fatalf("expected ErrCharacterWithoutHP, got %v", err)
		}
	})
}
