package campaign

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	app "questmaster-core/internal/app/campaign"
	domain "questmaster-core/internal/domain/campaign"
	"questmaster-core/internal/domain/rpg"
)

type CampaignRepositoryPG struct {
	db *pgxpool.Pool
}

func NewCampaignRepositoryPG(db *pgxpool.Pool) *CampaignRepositoryPG {
	return &CampaignRepositoryPG{db: db}
}

func (r *CampaignRepositoryPG) GetByDmId(
	userId rpg.UserID,
) ([]domain.Campaign, error) {
	rows, err := r.db.Query(context.Background(), `
        SELECT c.*
        FROM campaign c
        WHERE c.dm_id = $1
    `, userId.UUID())
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	record, err := pgx.CollectRows(rows, pgx.RowToStructByName[CampaignRow])
	if err != nil {
		return nil, err
	}

	domain := make([]domain.Campaign, 0)

	for _, c := range record {
		val, err := MapRowToDomain(c)
		if err != nil {
			return nil, err
		}
		domain = append(domain, val)
	}

	return domain, nil
}

func (r *CampaignRepositoryPG) GetByPlayerId(
	userId rpg.UserID,
) ([]domain.Campaign, error) {
	rows, err := r.db.Query(context.Background(), `
		SELECT DISTINCT c.*
		FROM campaign c
		JOIN character_sheet cs ON cs.campaign_id = c.id
		WHERE cs.player_id = $1
	`, userId.UUID())
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	record, err := pgx.CollectRows(rows, pgx.RowToStructByName[CampaignRow])
	if err != nil {
		return nil, err
	}

	domain := make([]domain.Campaign, 0)

	for _, c := range record {
		val, err := MapRowToDomain(c)
		if err != nil {
			return nil, err
		}
		domain = append(domain, val)
	}

	return domain, nil
}

func (r *CampaignRepositoryPG) FindBySlug(
	slug rpg.Slug,
) (*domain.Campaign, error) {
	rows, err := r.db.Query(context.Background(), `
		SELECT c.*
		FROM campaign c
		WHERE c.slug = $1
	`, slug.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	record, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[CampaignRow])
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

func (r *CampaignRepositoryPG) FindById(
	id domain.CampaignID,
) (*domain.Campaign, error) {
	rows, err := r.db.Query(context.Background(), `
		SELECT c.*
		FROM campaign c
		WHERE c.id = $1
	`, int(id))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	record, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[CampaignRow])
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

func (r *CampaignRepositoryPG) Create(input app.CreateCampaignInput) (domain.Campaign, error) {
	var overview any
	if input.Overview != nil {
		overview = input.Overview.String()
	} else {
		overview = nil
	}
	rows, err := r.db.Query(context.Background(), `
		INSERT INTO campaign(name, dm_id, trpg_system, overview) 
		VALUES($1, $2, $3, $4) RETURNING *
	`, input.Name.String(), input.Dm.UUID(), input.System, overview)
	if err != nil {
		return domain.Campaign{}, err
	}
	defer rows.Close()
	record, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[CampaignRow])
	if err != nil {
		return domain.Campaign{}, err
	}

	val, err := MapRowToDomain(record)
	if err != nil {
		return domain.Campaign{}, err
	}

	return val, nil
}

func (r *CampaignRepositoryPG) UpdateStatus(newStatus domain.CampaignStatus, campaignId domain.CampaignID) (domain.Campaign, error) {
	rows, err := r.db.Query(context.Background(), `
		UPDATE campaign SET status = $1 
		WHERE id = $2
		RETURNING *
	`, newStatus, int(campaignId))
	if err != nil {
		return domain.Campaign{}, err
	}
	defer rows.Close()
	record, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[CampaignRow])
	if err != nil {
		return domain.Campaign{}, err
	}

	val, err := MapRowToDomain(record)
	if err != nil {
		return domain.Campaign{}, err
	}

	return val, nil
}
