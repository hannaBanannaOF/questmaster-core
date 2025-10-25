package vo

import enum "questmaster-core/domain/enumerations"

type CreateCharacterSheet struct {
	CharacterName string          `json:"characterName"`
	TrpgSystem    enum.TrpgSystem `json:"system"`
	MaxHp         int             `json:"maxHp"`
}
