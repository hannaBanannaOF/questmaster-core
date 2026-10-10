package character

import (
	"context"
	"errors"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	campaignDomain "questmaster-core/internal/campaign/domain"
	characterDomain "questmaster-core/internal/character/domain"
	rpgDomain "questmaster-core/internal/rpg/domain"
	"questmaster-core/internal/shared/pagination"
	userDomain "questmaster-core/internal/user/domain"
)

type CharacterRepositoryPG struct {
	db *pgxpool.Pool
}

func NewCharacterRepositoryPG(db *pgxpool.Pool) *CharacterRepositoryPG {
	return &CharacterRepositoryPG{db: db}
}

// GetAllByPlayerIDWithFilters returns a page of the player characters matching the filters,
// ordered by name ignoring case and accents, then id, and the number of characters matching the filters.
func (r *CharacterRepositoryPG) GetAllByPlayerIDWithFilters(
	userID userDomain.UserID,
	filters characterDomain.CharacterListFilters,
	page pagination.Page,
) (pagination.Result[characterDomain.Character], error) {
	ctx := context.Background()
	where := " WHERE cs.player_id = $1"
	args := []any{userID.Value()}

	if filters.GameSystem != nil {
		args = append(args, filters.GameSystem.Value())
		where += " AND cs.game_system = $" + strconv.Itoa(len(args))
	}

	if filters.WithoutCampaign != nil {
		if *filters.WithoutCampaign {
			where += " AND cs.campaign_id IS NULL"
		} else {
			where += " AND cs.campaign_id IS NOT NULL"
		}
	}

	var total int
	if err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM character_sheet cs"+where, args...).Scan(&total); err != nil {
		return pagination.Result[characterDomain.Character]{}, err
	}

	args = append(args, page.Limit, page.Offset)
	rows, err := r.db.Query(ctx, "SELECT cs.* FROM character_sheet cs"+where+
		" ORDER BY unaccent(lower(cs.name)), cs.id"+
		" LIMIT $"+strconv.Itoa(len(args)-1)+" OFFSET $"+strconv.Itoa(len(args)), args...)
	if err != nil {
		return pagination.Result[characterDomain.Character]{}, err
	}
	defer rows.Close()

	record, err := pgx.CollectRows(rows, pgx.RowToStructByName[CharacterRow])
	if err != nil {
		return pagination.Result[characterDomain.Character]{}, err
	}

	items := make([]characterDomain.Character, 0, len(record))
	for _, c := range record {
		val, err := MapRowToDomain(c)
		if err != nil {
			return pagination.Result[characterDomain.Character]{}, err
		}
		items = append(items, val)
	}

	return pagination.Result[characterDomain.Character]{Items: items, Total: total}, nil
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

func (r *CharacterRepositoryPG) Create(name characterDomain.CharacterName, playerID userDomain.UserID, system rpgDomain.System, hp *characterDomain.HP) (characterDomain.Character, error) {

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
		INSERT INTO character_sheet(name, player_id, game_system, max_hp, current_hp) 
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

func (r *CharacterRepositoryPG) GetAllByUserIDAndCampaignIDNullAndSystem(userID userDomain.UserID, system rpgDomain.System) ([]characterDomain.Character, error) {
	rows, err := r.db.Query(context.Background(), `
        SELECT cs.*
        FROM character_sheet cs
        WHERE cs.campaign_id IS NULL
		AND cs.player_id = $1
		AND cs.game_system = $2
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

func (r *CharacterRepositoryPG) UpdateCampaign(campaignID campaignDomain.CampaignID, characterSlug rpgDomain.Slug, userID userDomain.UserID) (*characterDomain.Character, error) {
	rows, err := r.db.Query(context.Background(), `
		UPDATE character_sheet cs
		SET campaign_id = c.id
		FROM campaign c
		WHERE 
			cs.slug = $2
			AND cs.campaign_id IS NULL
			AND c.id = $1
			AND cs.player_id = $3
			AND cs.game_system = c.game_system
		RETURNING cs.*
	`, campaignID.Value(), characterSlug.Value(), userID.Value())
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
