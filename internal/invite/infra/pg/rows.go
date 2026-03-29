package invite

import "github.com/google/uuid"

type InviteRow struct {
	Id         int       `db:"id"`
	CampaignId int       `db:"campaign_id"`
	Hash       uuid.UUID `db:"hash"`
}
