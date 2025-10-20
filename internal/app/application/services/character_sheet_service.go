package services

import (
	"questmaster-core/domain/vo"
)

type CharacterSheetServiceInterface interface {
	GetAllByPlayerId(UserId string) []vo.CharacterSheetListItem
}
