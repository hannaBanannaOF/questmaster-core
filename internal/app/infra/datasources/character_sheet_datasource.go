package datasource

import (
	enum "questmaster-core/domain/enumerations"
	models "questmaster-core/internal/app/infra/models"
)

type CharacterSheetDataSourceInterface interface {
	GetAllByPlayerId(UserId string) ([]models.CharacterSheet, error)
	GetAllBySessionId(SessionId int) ([]models.CharacterSheet, error)
	GetOne(CharacterSheetId int) (*models.CharacterSheet, error)
	ResolveSlug(Slug string) (*int, error)
	CreateCharacterSheet(CharacterName string, MaxHp int, System enum.TrpgSystem, UserId string) (*models.CharacterSheet, error)
}
