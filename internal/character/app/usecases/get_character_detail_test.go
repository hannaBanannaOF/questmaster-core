package character

import (
	"errors"
	"testing"

	characterApp "questmaster-core/internal/character/app"
	characterDomain "questmaster-core/internal/character/domain"
	userDomain "questmaster-core/internal/user/domain"

	"github.com/google/uuid"
)

type fakeDetailRepo struct {
	characterApp.CharacterRepository
	character *characterDomain.Character
	recorded  []userDomain.UserID
	err       error
}

func (f *fakeDetailRepo) FindByID(characterDomain.CharacterID) (*characterDomain.Character, error) {
	return f.character, nil
}

func (f *fakeDetailRepo) RecordAccess(user userDomain.UserID, _ characterDomain.CharacterID) error {
	f.recorded = append(f.recorded, user)
	return f.err
}

func TestGetCharacterDetailRecordsOwnOpens(t *testing.T) {
	player := userDomain.NewUserID(uuid.New())
	newRepo := func() *fakeDetailRepo {
		return &fakeDetailRepo{character: &characterDomain.Character{Id: 1, Name: "Hero", PlayerID: player}}
	}

	t.Run("owner opening is recorded", func(t *testing.T) {
		repo := newRepo()
		if _, err := NewGetCharacterDetail(repo).Execute(characterApp.GetCharacterDetailsCommand{ID: 1, UserID: player}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(repo.recorded) != 1 || repo.recorded[0] != player {
			t.Fatalf("expected one access recorded for the owner, got %v", repo.recorded)
		}
	})

	t.Run("another user is not recorded", func(t *testing.T) {
		repo := newRepo()
		other := userDomain.NewUserID(uuid.New())
		if _, err := NewGetCharacterDetail(repo).Execute(characterApp.GetCharacterDetailsCommand{ID: 1, UserID: other}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(repo.recorded) != 0 {
			t.Fatalf("expected nothing recorded, got %v", repo.recorded)
		}
	})

	t.Run("failing recorder still returns the character", func(t *testing.T) {
		repo := newRepo()
		repo.err = errors.New("database is down")
		got, err := NewGetCharacterDetail(repo).Execute(characterApp.GetCharacterDetailsCommand{ID: 1, UserID: player})
		if err != nil || got.Id != 1 {
			t.Fatalf("expected the character despite the recorder error, got %+v err=%v", got, err)
		}
	})
}
