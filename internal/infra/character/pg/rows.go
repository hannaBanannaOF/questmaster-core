package character

import (
	"github.com/google/uuid"
)

type CharacterRow struct {
	Id         int       `db:"id"`
	Name       string    `db:"name"`
	PlayerId   uuid.UUID `db:"player_id"`
	System     string    `db:"trpg_system"`
	CampaingId *int      `db:"campaign_id"`
	Slug       string    `db:"slug"`
	MaxHp      *int      `db:"max_hp"`
	CurrentHp  *int      `db:"current_hp"`
}
