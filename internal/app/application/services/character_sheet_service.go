package services

import (
	enum "questmaster-core/domain/enumerations"
	"questmaster-core/domain/vo"
)

type CharacterSheetServiceInterface interface {
	GetAllByPlayerId(UserId string) []vo.CharacterSheetListItem
	GetCharacterSheetDetails(CharacterSheetId int) *vo.CharacterSheetDetailItem
	ResolveSlug(Slug string) *vo.SlugResolve
	CreateCharacterSheet(CharacterName string, MaxHp int, System enum.TrpgSystem, UserId string) *vo.Slug
}
