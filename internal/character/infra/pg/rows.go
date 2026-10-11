package character

import (
	"time"

	"github.com/google/uuid"
)

type CharacterRow struct {
	Id         int       `db:"id"`
	Name       string    `db:"name"`
	PlayerID   uuid.UUID `db:"player_id"`
	System     string    `db:"game_system"`
	CampaingID *int      `db:"campaign_id"`
	Slug       string    `db:"slug"`
	MaxHp      *int      `db:"max_hp"`
	CurrentHp  *int      `db:"current_hp"`
}

// CharacterListRow is a character of the requester list, with their last access to it
type CharacterListRow struct {
	CharacterRow
	LastAccessedAt *time.Time `db:"last_accessed_at"`
}
