package models

import (
	enum "questmaster-core/domain/enumerations"

	"github.com/google/uuid"
)

type Session struct {
	Id          int             `db:"id"`
	SessionName string          `db:"session_name"`
	DmId        uuid.UUID       `db:"dm_id"`
	InPlay      bool            `db:"in_play"`
	TrpgSystem  enum.TrpgSystem `db:"trpg_system"`
	Slug        string          `db:"slug"`
}
