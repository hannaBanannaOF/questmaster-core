package campaign

import (
	"github.com/google/uuid"
)

type CampaignRow struct {
	ID          int       `db:"id"`
	Name        string    `db:"name"`
	DmID        uuid.UUID `db:"dm_id"`
	Status      string    `db:"status"`
	System      string    `db:"game_system"`
	Slug        string    `db:"slug"`
	Overview    *string   `db:"overview"`
	PlayerCount int       `db:"player_count"`
}
