package vo

import (
	enum "questmaster-core/domain/enumerations"
)

type CharacterSheetListItem struct {
	Slug        string          `json:"slug"`
	Description string          `json:"description"`
	System      enum.TrpgSystem `json:"system"`
}
