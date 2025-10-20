package services_v1

import (
	"log"
	vo "questmaster-core/domain/vo"
	datasource "questmaster-core/internal/app/infra/datasources"
	models "questmaster-core/internal/app/infra/models"

	"github.com/samber/lo"
)

type CharacterSheetServiceV1 struct {
	CharacterSheetDs datasource.CharacterSheetDataSourceInterface
}

func (svc *CharacterSheetServiceV1) GetAllByPlayerId(UserId string) []vo.CharacterSheetListItem {
	data, err := svc.CharacterSheetDs.GetAllByPlayerId(UserId)
	if err != nil {
		log.Panicf("Unable to get character sheets: %s", err)
	}

	return lo.Map(data, func(model models.CharacterSheet, _ int) vo.CharacterSheetListItem {
		return vo.CharacterSheetListItem{
			Slug:        model.Slug,
			Description: model.CharacterName,
			System:      model.TrpgSystem,
		}
	})
}
