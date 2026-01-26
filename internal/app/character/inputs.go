package character

import (
	domain "questmaster-core/internal/domain/character"
	"questmaster-core/internal/domain/rpg"
)

type CreateCharacterInput struct {
	Name   domain.CharacterName
	Hp     domain.HP
	System rpg.System
	Player rpg.UserID
}
