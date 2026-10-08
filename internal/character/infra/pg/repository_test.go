package character

import (
	"regexp"
	"testing"

	characterDomain "questmaster-core/internal/character/domain"
	rpgDomain "questmaster-core/internal/rpg/domain"
	"questmaster-core/internal/shared/testdb"
	userDomain "questmaster-core/internal/user/domain"

	"github.com/google/uuid"
)

func createCharacter(t *testing.T, r *CharacterRepositoryPG, name string, player userDomain.UserID) characterDomain.Character {
	t.Helper()
	characterName, err := characterDomain.NewCharacterName(name)
	if err != nil {
		t.Fatalf("character name: %v", err)
	}
	c, err := r.Create(characterName, player, rpgDomain.DungeonsAndDragons, nil)
	if err != nil {
		t.Fatalf("create character: %v", err)
	}
	return c
}

func TestCharacterSlugs(t *testing.T) {
	r := NewCharacterRepositoryPG(testdb.Pool(t))

	t.Run("duplicate names get a numeric suffix", func(t *testing.T) {
		// Random suffix keeps the name unique across test runs on the same database
		name := "Twin " + uuid.NewString()[:8]
		player := userDomain.NewUserID(uuid.New())
		first := createCharacter(t, r, name, player)
		second := createCharacter(t, r, name, player)
		if second.Slug.Value() != first.Slug.Value()+"-2" {
			t.Fatalf("expected %q, got %q", first.Slug.Value()+"-2", second.Slug.Value())
		}
	})

	t.Run("fallback for names without ASCII letters or digits", func(t *testing.T) {
		player := userDomain.NewUserID(uuid.New())
		c := createCharacter(t, r, "東京", player)
		if !regexp.MustCompile(`^character(-\d+)?$`).MatchString(c.Slug.Value()) {
			t.Fatalf("unexpected slug %q", c.Slug.Value())
		}

		list, err := r.GetAllByPlayerIDWithFilters(player, &characterDomain.CharacterListFilters{})
		if err != nil || len(list) != 1 || list[0].Id != c.Id {
			t.Fatalf("expected the character in the player list, got %v err=%v", list, err)
		}
	})
}
