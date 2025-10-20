package vo

import (
	enum "questmaster-core/domain/enumerations"
)

type SessionListItem struct {
	Slug        string          `json:"slug"`
	Description string          `json:"description"`
	Dmed        bool            `json:"dmed"`
	System      enum.TrpgSystem `json:"system"`
}
