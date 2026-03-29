package character

import (
	characterDomain "questmaster-core/internal/character/domain"
	rpgDomain "questmaster-core/internal/rpg/domain"
)

type CreateCharacterCommand struct {
	Name   characterDomain.CharacterName
	Hp     *characterDomain.HP
	System rpgDomain.System
	Player rpgDomain.UserID
}

type DeleteCharacterCommand struct {
	ID     characterDomain.CharacterID
	UserID rpgDomain.UserID
}

type UpdateHPCommand struct {
	ID     characterDomain.CharacterID
	NewHP  characterDomain.HP
	UserID rpgDomain.UserID
}
