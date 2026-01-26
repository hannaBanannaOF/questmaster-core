package character

import (
	domain "questmaster-core/internal/domain/character"
	"questmaster-core/internal/domain/rpg"
)

type CharacterRepository interface {
	GetAllByPlayerId(userId rpg.UserID) ([]domain.Character, error)
	FindBySlug(slug rpg.Slug) (*domain.Character, error)
	FindById(characterId domain.CharacterID) (*domain.Character, error)
	Create(input CreateCharacterInput) (domain.Character, error)
}
