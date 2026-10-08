package character

import (
	"errors"
	"testing"

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

func TestUpdateHPWithoutHP(t *testing.T) {
	player := userDomain.NewUserID(uuid.New())
	repo := fakeCharacterRepo{character: &characterDomain.Character{Id: 1, PlayerID: player}}

	_, err := NewUpdateHP(repo, nil).Execute(characterApp.UpdateHPCommand{ID: 1, NewHP: 5, UserID: player})
	if !errors.Is(err, characterDomain.ErrCharacterWithoutHP) {
		t.Fatalf("expected ErrCharacterWithoutHP, got %v", err)
	}
}
