package character

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	campaignDomain "questmaster-core/internal/campaign/domain"
	characterDomain "questmaster-core/internal/character/domain"
	rpgDomain "questmaster-core/internal/rpg/domain"
)

type CharacterRepositoryPG struct {
	db *pgxpool.Pool
}

func NewCharacterRepositoryPG(db *pgxpool.Pool) *CharacterRepositoryPG {
	return &CharacterRepositoryPG{db: db}
}

func (r *CharacterRepositoryPG) GetAllByPlayerID(
	userID rpgDomain.UserID,
) ([]characterDomain.Character, error) {
	rows, err := r.db.Query(context.Background(), `
        SELECT cs.*
        FROM character_sheet cs
        WHERE cs.player_id = $1
    `, userID.Value())
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	record, err := pgx.CollectRows(rows, pgx.RowToStructByName[CharacterRow])
	if err != nil {
		return nil, err
	}

	domain := make([]characterDomain.Character, 0)

	for _, c := range record {
		val, err := MapRowToDomain(c)
		if err != nil {
			return nil, err
		}
		domain = append(domain, val)
	}

	return domain, nil
}

func (r *CharacterRepositoryPG) GetAllByCampaignID(campaignID campaignDomain.CampaignID) ([]characterDomain.Character, error) {
	rows, err := r.db.Query(context.Background(), `
        SELECT cs.*
        FROM character_sheet cs
        WHERE cs.campaign_id = $1
    `, campaignID.Value())
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	record, err := pgx.CollectRows(rows, pgx.RowToStructByName[CharacterRow])
	if err != nil {
		return nil, err
	}

	domain := make([]characterDomain.Character, 0)

	for _, c := range record {
		val, err := MapRowToDomain(c)
		if err != nil {
			return nil, err
		}
		domain = append(domain, val)
	}

	return domain, nil
}

func (r *CharacterRepositoryPG) FindBySlug(slug rpgDomain.Slug) (*characterDomain.Character, error) {
	rows, err := r.db.Query(context.Background(), `
		SELECT cs.*
		FROM character_sheet cs
		WHERE cs.slug = $1
	`, slug.Value())
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	record, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[CharacterRow])
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

func (r *CharacterRepositoryPG) FindByID(characterID characterDomain.CharacterID) (*characterDomain.Character, error) {
	rows, err := r.db.Query(context.Background(), `
		SELECT cs.*
		FROM character_sheet cs
		WHERE cs.id = $1
	`, characterID.Value())
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	record, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[CharacterRow])
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

func (r *CharacterRepositoryPG) Create(name characterDomain.CharacterName, playerID rpgDomain.UserID, system rpgDomain.System, hp *characterDomain.HP) (characterDomain.Character, error) {

	var currHP any
	var maxHP any
	if hp != nil {
		currHP = hp.Current()
		maxHP = hp.Max()
	} else {
		currHP = nil
		maxHP = nil
	}

	rows, err := r.db.Query(context.Background(), `
		INSERT INTO character_sheet(name, player_id, trpg_system, max_hp, current_hp) 
		VALUES($1, $2, $3, $4, $5) 
		RETURNING *
	`, name.Value(), playerID.Value(), system.Value(), maxHP, currHP)
	if err != nil {
		return characterDomain.Character{}, err
	}
	defer rows.Close()
	record, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[CharacterRow])
	if err != nil {
		return characterDomain.Character{}, err
	}

	val, err := MapRowToDomain(record)
	if err != nil {
		return characterDomain.Character{}, err
	}

	return val, nil
}

func (r *CharacterRepositoryPG) UpdateHP(newHP characterDomain.HP, characterID characterDomain.CharacterID) (characterDomain.Character, error) {
	rows, err := r.db.Query(context.Background(), `
		UPDATE character_sheet SET current_hp = $1, max_hp = $2
		WHERE id = $3
		RETURNING *
	`, newHP.Current(), newHP.Max(), characterID.Value())
	if err != nil {
		return characterDomain.Character{}, err
	}
	defer rows.Close()
	record, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[CharacterRow])
	if err != nil {
		return characterDomain.Character{}, err
	}

	val, err := MapRowToDomain(record)
	if err != nil {
		return characterDomain.Character{}, err
	}

	return val, nil
}

func (r *CharacterRepositoryPG) DeleteByID(characterID characterDomain.CharacterID) (bool, error) {
	cmdTag, err := r.db.Exec(context.Background(), `
		DELETE FROM character_sheet
		WHERE id = $1
	`, characterID.Value())
	if err != nil {
		return false, err
	}

	return cmdTag.RowsAffected() > 0, nil
}

func (r *CharacterRepositoryPG) GetAllByUserIDAndCampaignIDNullAndSystem(userID rpgDomain.UserID, system rpgDomain.System) ([]characterDomain.Character, error) {
	rows, err := r.db.Query(context.Background(), `
        SELECT cs.*
        FROM character_sheet cs
        WHERE cs.campaign_id IS NULL
		AND cs.player_id = $1
		AND cs.trpg_system = $2
    `, userID.Value(), system.Value())
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	record, err := pgx.CollectRows(rows, pgx.RowToStructByName[CharacterRow])
	if err != nil {
		return nil, err
	}

	domain := make([]characterDomain.Character, 0)

	for _, c := range record {
		val, err := MapRowToDomain(c)
		if err != nil {
			return nil, err
		}
		domain = append(domain, val)
	}

	return domain, nil
}

func (r *CharacterRepositoryPG) UpdateCampaign(campaignID campaignDomain.CampaignID, characterID characterDomain.CharacterID) (*characterDomain.Character, error) {
	rows, err := r.db.Query(context.Background(), `
		UPDATE character_sheet cs
		SET campaign_id = $1
		FROM campaign c
		WHERE 
			cs.id = $2
			AND cs.campaign_id IS NULL
			AND c.id = $1
			AND cs.trpg_system = c.trpg_system
		RETURNING cs.*
	`, campaignID.Value(), characterID.Value())
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	record, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[CharacterRow])
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
