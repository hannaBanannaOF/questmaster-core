package character

import (
	"time"

	"github.com/google/uuid"
)

type CharacterRow struct {
	Id         int           `db:"id"`
	Name       string        `db:"name"`
	PlayerID   uuid.NullUUID `db:"player_id"`
	PlayerName *string       `db:"player_name"`
	System     string        `db:"game_system"`
	CampaingID *int          `db:"campaign_id"`
	Slug       string        `db:"slug"`
}

// CharacterListRow is a character of the requester list, with their last access to it
type CharacterListRow struct {
	CharacterRow
	LastAccessedAt *time.Time `db:"last_accessed_at"`
	CampaignSlug   *string    `db:"campaign_slug"`
	CampaignName   *string    `db:"campaign_name"`
}
