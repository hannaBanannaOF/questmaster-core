package campaign

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	campaignDomain "questmaster-core/internal/campaign/domain"
	rpgDomain "questmaster-core/internal/rpg/domain"
	userDomain "questmaster-core/internal/user/domain"
)

type CampaignRepositoryPG struct {
	db *pgxpool.Pool
}

func NewCampaignRepositoryPG(db *pgxpool.Pool) *CampaignRepositoryPG {
	return &CampaignRepositoryPG{db: db}
}

func (r *CampaignRepositoryPG) GetByDmId(userID userDomain.UserID) ([]campaignDomain.Campaign, error) {
	rows, err := r.db.Query(context.Background(), `
        SELECT c.*, COUNT(cs.id) as player_count
        FROM campaign c
		LEFT JOIN character_sheet cs ON cs.campaign_id = c.id
		WHERE c.dm_id = $1
		GROUP BY c.id
    `, userID.Value())
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	record, err := pgx.CollectRows(rows, pgx.RowToStructByName[CampaignRow])
	if err != nil {
		return nil, err
	}

	domain := make([]campaignDomain.Campaign, 0)

	for _, c := range record {
		val, err := MapRowToDomain(c)
		if err != nil {
			return nil, err
		}
		domain = append(domain, val)
	}

	return domain, nil
}

func (r *CampaignRepositoryPG) GetByPlayerId(userID userDomain.UserID) ([]campaignDomain.Campaign, error) {
	rows, err := r.db.Query(context.Background(), `
		SELECT DISTINCT c.*, COUNT(cs.id) as player_count
        FROM campaign c
		LEFT JOIN character_sheet cs ON cs.campaign_id = c.id
		WHERE cs.player_id = $1
		GROUP BY c.id
	`, userID.Value())
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	record, err := pgx.CollectRows(rows, pgx.RowToStructByName[CampaignRow])
	if err != nil {
		return nil, err
	}

	domain := make([]campaignDomain.Campaign, 0)

	for _, c := range record {
		val, err := MapRowToDomain(c)
		if err != nil {
			return nil, err
		}
		domain = append(domain, val)
	}

	return domain, nil
}

func (r *CampaignRepositoryPG) FindBySlug(slug rpgDomain.Slug) (*campaignDomain.Campaign, error) {
	rows, err := r.db.Query(context.Background(), `
		SELECT c.*, COUNT(cs.id) as player_count
        FROM campaign c
		LEFT JOIN character_sheet cs ON cs.campaign_id = c.id
		WHERE c.slug = $1
		GROUP BY c.id	
	`, slug.Value())
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

func (r *CampaignRepositoryPG) FindById(id campaignDomain.CampaignID) (*campaignDomain.Campaign, error) {
	rows, err := r.db.Query(context.Background(), `
		SELECT c.*, COUNT(cs.id) as player_count
        FROM campaign c
		LEFT JOIN character_sheet cs ON cs.campaign_id = c.id
		WHERE c.id = $1
		GROUP BY c.id
	`, id.Value())
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

func (r *CampaignRepositoryPG) Create(Name campaignDomain.CampaignName, Overview *campaignDomain.CampaignOverview, DmID userDomain.UserID, System rpgDomain.System) (campaignDomain.Campaign, error) {
	var overview any
	if Overview != nil {
		overview = Overview.Value()
	} else {
		overview = nil
	}
	rows, err := r.db.Query(context.Background(), `
		INSERT INTO campaign(name, dm_id, game_system, overview) 
		VALUES($1, $2, $3, $4) RETURNING *, 0 as player_count
	`, Name.Value(), DmID.Value(), System.Value(), overview)
	if err != nil {
		return campaignDomain.Campaign{}, err
	}
	defer rows.Close()
	record, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[CampaignRow])
	if err != nil {
		return campaignDomain.Campaign{}, err
	}

	val, err := MapRowToDomain(record)
	if err != nil {
		return campaignDomain.Campaign{}, err
	}

	return val, nil
}

func (r *CampaignRepositoryPG) UpdateStatus(newStatus campaignDomain.CampaignStatus, id campaignDomain.CampaignID) (campaignDomain.Campaign, error) {
	rows, err := r.db.Query(context.Background(), `
		UPDATE campaign SET status = $1 
		WHERE id = $2
		RETURNING *, 0 as player_count
	`, newStatus.Value(), id.Value())
	if err != nil {
		return campaignDomain.Campaign{}, err
	}
	defer rows.Close()
	record, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[CampaignRow])
	if err != nil {
		return campaignDomain.Campaign{}, err
	}

	val, err := MapRowToDomain(record)
	if err != nil {
		return campaignDomain.Campaign{}, err
	}

	return val, nil
}

func (r *CampaignRepositoryPG) DeleteById(id campaignDomain.CampaignID) (bool, error) {
	cmdTag, err := r.db.Exec(context.Background(), `
		DELETE FROM campaign
		WHERE id = $1
	`, id.Value())
	if err != nil {
		return false, err
	}

	return cmdTag.RowsAffected() > 0, nil
}
