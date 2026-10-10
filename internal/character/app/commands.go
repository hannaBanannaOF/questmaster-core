package character

import (
	characterDomain "questmaster-core/internal/character/domain"
	rpgDomain "questmaster-core/internal/rpg/domain"
	"questmaster-core/internal/shared/pagination"
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
	UserID  userDomain.UserID
	Filters characterDomain.CharacterListFilters
	Page    pagination.Page
}

type UpdateHPCommand struct {
	ID     characterDomain.CharacterID
	NewHP  int
	UserID userDomain.UserID
}
