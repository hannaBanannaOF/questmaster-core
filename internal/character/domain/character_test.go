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

func TestEdit(t *testing.T) {
	player, dm, other := userDomain.NewUserID(uuid.New()), userDomain.NewUserID(uuid.New()), userDomain.NewUserID(uuid.New())
	campaign := fakeCampaign{dm: dm}
	ptr := func(v int) *int { return &v }
	name := CharacterName("Harvey Walters")
	withHP := func(current, max int) Character {
		hp, _ := NewHP(current, max)
		return Character{PlayerID: &player, Hp: &hp}
	}

	t.Run("player renames", func(t *testing.T) {
		c := withHP(10, 12)
		if err := c.Edit(player, nil, &name, nil); err != nil || c.Name != name || c.Hp.Current() != 10 {
			t.Fatalf("expected renamed with HP kept, got %+v err=%v", c, err)
		}
	})

	t.Run("campaign DM raises max HP keeping current", func(t *testing.T) {
		c := withHP(10, 10)
		if err := c.Edit(dm, campaign, nil, ptr(13)); err != nil || c.Hp.Max() != 13 || c.Hp.Current() != 10 {
			t.Fatalf("expected 10 / 13, got %+v err=%v", c.Hp, err)
		}
	})

	t.Run("max below current lowers current", func(t *testing.T) {
		c := withHP(10, 12)
		if err := c.Edit(player, nil, nil, ptr(8)); err != nil || c.Hp.Max() != 8 || c.Hp.Current() != 8 {
			t.Fatalf("expected 8 / 8, got %+v err=%v", c.Hp, err)
		}
	})

	t.Run("character without HP gets full HP", func(t *testing.T) {
		c := Character{PlayerID: &player}
		if err := c.Edit(player, nil, nil, ptr(12)); err != nil || c.Hp == nil || c.Hp.Current() != 12 || c.Hp.Max() != 12 {
			t.Fatalf("expected 12 / 12, got %+v err=%v", c.Hp, err)
		}
	})

	t.Run("DM edits an unclaimed premade character", func(t *testing.T) {
		c := Character{}
		if err := c.Edit(dm, campaign, &name, nil); err != nil || c.Name != name {
			t.Fatalf("expected the DM to rename the premade, got err=%v", err)
		}
	})

	t.Run("other user is forbidden and nothing changes", func(t *testing.T) {
		c := withHP(10, 12)
		if err := c.Edit(other, campaign, &name, ptr(1)); !errors.Is(err, ErrNotAllowed) || c.Name == name || c.Hp.Max() != 12 {
			t.Fatalf("expected ErrNotAllowed with no change, got %+v err=%v", c, err)
		}
	})

	t.Run("non-positive max HP is rejected and nothing changes", func(t *testing.T) {
		c := withHP(10, 12)
		if err := c.Edit(player, nil, &name, ptr(0)); !errors.Is(err, ErrInvalidMaxHP) || c.Name == name || c.Hp.Max() != 12 {
			t.Fatalf("expected ErrInvalidMaxHP with no change, got %+v err=%v", c, err)
		}
	})
}
