package services_v1

import (
	"log"
	enum "questmaster-core/domain/enumerations"
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

func (svc *CharacterSheetServiceV1) GetCharacterSheetDetails(CharacterSheetId int) *vo.CharacterSheetDetailItem {
	data, err := svc.CharacterSheetDs.GetOne(CharacterSheetId)
	if err != nil {
		log.Panicf("Unable to get character sheet: %s", err)
	}

	if data == nil {
		return nil
	}

	return &vo.CharacterSheetDetailItem{
		Id:        data.Id,
		Name:      data.CharacterName,
		System:    data.TrpgSystem,
		MaxHP:     data.MaxHp,
		CurrentHP: data.CurrentHp,
	}
}

func (svc *CharacterSheetServiceV1) ResolveSlug(Slug string) *vo.SlugResolve {
	data, err := svc.CharacterSheetDs.ResolveSlug(Slug)
	if err != nil {
		log.Panicf("Unable to resolve slug %s: %s", Slug, err)
	}

	return &vo.SlugResolve{
		CoreId: *data,
	}
}

func (svc *CharacterSheetServiceV1) CreateCharacterSheet(CharacterName string, MaxHp int, System enum.TrpgSystem, UserId string) *vo.Slug {
	data, err := svc.CharacterSheetDs.CreateCharacterSheet(CharacterName, MaxHp, System, UserId)
	if err != nil {
		log.Panicf("Unable to create character sheet: %s", err)
	}
	return &vo.Slug{
		Slug: data.Slug,
	}
}
