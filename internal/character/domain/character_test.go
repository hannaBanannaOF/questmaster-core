package character

import (
	"errors"
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

func TestRename(t *testing.T) {
	player, dm, other := userDomain.NewUserID(uuid.New()), userDomain.NewUserID(uuid.New()), userDomain.NewUserID(uuid.New())
	campaign := fakeCampaign{dm: dm}
	name := CharacterName("Harvey Walters")

	t.Run("player renames", func(t *testing.T) {
		c := Character{PlayerID: &player}
		if err := c.Rename(player, nil, name); err != nil || c.Name != name {
			t.Fatalf("expected renamed, got %+v err=%v", c, err)
		}
	})

	t.Run("DM renames an unclaimed premade character", func(t *testing.T) {
		c := Character{}
		if err := c.Rename(dm, campaign, name); err != nil || c.Name != name {
			t.Fatalf("expected the DM to rename the premade, got err=%v", err)
		}
	})

	t.Run("other user is forbidden and nothing changes", func(t *testing.T) {
		c := Character{Name: "Harvey", PlayerID: &player}
		if err := c.Rename(other, campaign, name); !errors.Is(err, ErrNotAllowed) || c.Name == name {
			t.Fatalf("expected ErrNotAllowed with no change, got %+v err=%v", c, err)
		}
	})
}
