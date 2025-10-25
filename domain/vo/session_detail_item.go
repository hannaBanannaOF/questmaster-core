package vo

import enum "questmaster-core/domain/enumerations"

type SessionCharacterSheetItem struct {
	Name      string `json:"name"`
	MaxHp     *int   `json:"maxHp"`
	CurrentHp *int   `json:"currentHp"`
}

type SessionDetailItem struct {
	Id         int                         `json:"id"`
	Name       string                      `json:"name"`
	Overview   *string                     `json:"overview"`
	System     enum.TrpgSystem             `json:"system"`
	Dmed       bool                        `json:"dmed"`
	InPlay     bool                        `json:"inPlay"`
	Characters []SessionCharacterSheetItem `json:"characters"`
}
