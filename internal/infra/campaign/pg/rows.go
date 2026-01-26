package campaign

import (
	"github.com/google/uuid"
)

type CampaignRow struct {
	Id       int       `db:"id"`
	Name     string    `db:"name"`
	DmId     uuid.UUID `db:"dm_id"`
	Status   string    `db:"status"`
	System   string    `db:"trpg_system"`
	Slug     string    `db:"slug"`
	Overview *string   `db:"overview"`
}
