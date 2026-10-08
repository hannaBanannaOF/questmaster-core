package character

import "testing"

func TestMapRowToDomainHP(t *testing.T) {
	ten := 10
	base := CharacterRow{Id: 1, Name: "Hero", System: "DUNGEONS_AND_DRAGONS", Slug: "hero"}

	t.Run("no hp", func(t *testing.T) {
		c, err := MapRowToDomain(base)
		if err != nil || c.Hp != nil {
			t.Fatalf("expected character without hp, got hp=%v err=%v", c.Hp, err)
		}
	})

	t.Run("both hp columns", func(t *testing.T) {
		row := base
		row.CurrentHp, row.MaxHp = &ten, &ten
		c, err := MapRowToDomain(row)
		if err != nil || c.Hp == nil || c.Hp.Max() != 10 {
			t.Fatalf("expected hp 10/10, got hp=%v err=%v", c.Hp, err)
		}
	})

	t.Run("only one hp column returns error instead of panicking", func(t *testing.T) {
		row := base
		row.MaxHp = &ten
		if _, err := MapRowToDomain(row); err == nil {
			t.Fatalf("expected error for inconsistent hp")
		}
	})
}
