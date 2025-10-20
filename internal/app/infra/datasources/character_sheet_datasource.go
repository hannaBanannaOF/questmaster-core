package datasource

import models "questmaster-core/internal/app/infra/models"

type CharacterSheetDataSourceInterface interface {
	GetAllByPlayerId(UserId string) ([]models.CharacterSheet, error)
}
