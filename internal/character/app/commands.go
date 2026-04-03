package character

import (
	characterDomain "questmaster-core/internal/character/domain"
	rpgDomain "questmaster-core/internal/rpg/domain"
	userDomain "questmaster-core/internal/user/domain"
)

type CreateCharacterCommand struct {
	Name   characterDomain.CharacterName
	Hp     *characterDomain.HP
	System rpgDomain.System
	Player userDomain.UserID
}

type DeleteCharacterCommand struct {
	ID     characterDomain.CharacterID
	UserID userDomain.UserID
}

type GetCharacterDetailsCommand struct {
	ID characterDomain.CharacterID
}

type GetCurrentUserCharactersCommand struct {
	UserID userDomain.UserID
}

type UpdateHPCommand struct {
	ID     characterDomain.CharacterID
	NewHP  characterDomain.HP
	UserID userDomain.UserID
}
