package invite

import (
	"context"
	"errors"
	campaignDomain "questmaster-core/internal/campaign/domain"
	inviteDomain "questmaster-core/internal/invite/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type InviteRepositoryPG struct {
	db *pgxpool.Pool
}

func NewInviteRepositoryPG(db *pgxpool.Pool) *InviteRepositoryPG {
	return &InviteRepositoryPG{db: db}
}

func (r *InviteRepositoryPG) Create(campaignID campaignDomain.CampaignID) (*inviteDomain.Invite, error) {
	rows, err := r.db.Query(context.Background(), `
		INSERT INTO campaign_invite(campaign_id) 
		VALUES($1) ON CONFLICT DO NOTHING
		RETURNING *
	`, campaignID.Value())
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	record, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[InviteRow])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	val, err := MapRowToDomain(record)
	if err != nil {
		return nil, err
	}

	return &val, nil
}

func (r *InviteRepositoryPG) FindByCampaignID(campaignID campaignDomain.CampaignID) (*inviteDomain.Invite, error) {
	rows, err := r.db.Query(context.Background(), `
		SELECT ci.*
		FROM campaign_invite ci
		WHERE ci.campaign_id = $1
	`, campaignID.Value())
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	record, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[InviteRow])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	val, err := MapRowToDomain(record)
	if err != nil {
		return nil, err
	}

	return &val, nil
}

func (r *InviteRepositoryPG) FindByHash(hash inviteDomain.InviteHash) (*inviteDomain.Invite, error) {
	rows, err := r.db.Query(context.Background(), `
		SELECT ci.* 
		FROM campaign_invite ci
		WHERE ci.hash = $1
	`, hash.Value())
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	record, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[InviteRow])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	val, err := MapRowToDomain(record)
	if err != nil {
		return nil, err
	}

	return &val, nil
}
