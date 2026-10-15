package campaign

import (
	"time"

	"github.com/google/uuid"
)

type CampaignRow struct {
	ID          int       `db:"id"`
	Name        string    `db:"name"`
	DmID        uuid.UUID `db:"dm_id"`
	DmName      string    `db:"dm_name"`
	Status      string    `db:"status"`
	System      string    `db:"game_system"`
	Slug        string    `db:"slug"`
	Overview    *string   `db:"overview"`
	PlayerCount int       `db:"player_count"`
}

// CampaignListRow is a campaign of the requester list, with their last access to it
type CampaignListRow struct {
	CampaignRow
	LastAccessedAt *time.Time `db:"last_accessed_at"`
}
