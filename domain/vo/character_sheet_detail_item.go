package vo

import enum "questmaster-core/domain/enumerations"

type CharacterSheetDetailItem struct {
	Id        int             `json:"id"`
	Name      string          `json:"name"`
	System    enum.TrpgSystem `json:"system"`
	MaxHP     *int            `json:"maxHp"`
	CurrentHP *int            `json:"currentHp"`
}
