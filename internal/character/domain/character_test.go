package character

import (
	"testing"

	userDomain "questmaster-core/internal/user/domain"

	"github.com/google/uuid"
)

func TestIsPlayerAndIsPremade(t *testing.T) {
	player := userDomain.NewUserID(uuid.New())
	other := userDomain.NewUserID(uuid.New())

	owned := Character{PlayerID: &player}
	if !owned.IsPlayer(player) || owned.IsPlayer(other) || owned.IsPremade() {
		t.Fatalf("expected an owned character to belong only to its player")
	}

	premade := Character{}
	if premade.IsPlayer(player) || !premade.IsPremade() {
		t.Fatalf("expected a character without a player to be premade and nobody's")
	}
}
