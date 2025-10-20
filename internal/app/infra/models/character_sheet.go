package models

import (
	enum "questmaster-core/domain/enumerations"

	"github.com/google/uuid"
)

type CharacterSheet struct {
	Id            int             `db:"id"`
	CharacterName string          `db:"character_name"`
	PlayerId      uuid.UUID       `db:"player_id"`
	TrpgSystem    enum.TrpgSystem `db:"trpg_system"`
	SessionId     *int            `db:"session_id"`
	Slug          string          `db:"slug"`
}
