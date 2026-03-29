package character

import (
	"github.com/google/uuid"
)

type CharacterRow struct {
	Id         int       `db:"id"`
	Name       string    `db:"name"`
	PlayerID   uuid.UUID `db:"player_id"`
	System     string    `db:"trpg_system"`
	CampaingID *int      `db:"campaign_id"`
	Slug       string    `db:"slug"`
	MaxHp      *int      `db:"max_hp"`
	CurrentHp  *int      `db:"current_hp"`
}
