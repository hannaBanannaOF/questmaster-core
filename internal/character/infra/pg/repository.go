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

	if filters.Search != nil {
		args = append(args, *filters.Search)
		where += " AND " + pagination.NameContains("cs.name", len(args))
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
	// $1 is the requester, who is also the player filtered by where
	rows, err := r.db.Query(ctx, "SELECT cs.*, a.accessed_at AS last_accessed_at, lc.slug AS campaign_slug, lc.name AS campaign_name"+
		" FROM character_sheet cs LEFT JOIN character_access a ON a.user_id = $1 AND a.character_id = cs.id"+
		" LEFT JOIN campaign lc ON lc.id = cs.campaign_id"+where+
		page.Sort.OrderBy("a.accessed_at", "cs.name", "cs.id")+
		" LIMIT $"+strconv.Itoa(len(args)-1)+" OFFSET $"+strconv.Itoa(len(args)), args...)
	if err != nil {
		return pagination.Result[characterDomain.Character]{}, err
	}
	defer rows.Close()

	record, err := pgx.CollectRows(rows, pgx.RowToStructByName[CharacterListRow])
	if err != nil {
		return pagination.Result[characterDomain.Character]{}, err
	}

	items := make([]characterDomain.Character, 0, len(record))
	for _, c := range record {
		val, err := MapRowToDomain(c.CharacterRow)
		if err != nil {
			return pagination.Result[characterDomain.Character]{}, err
		}
		val.LastAccessedAt = c.LastAccessedAt
		if c.CampaignSlug != nil && c.CampaignName != nil {
			val.Campaign = &characterDomain.CampaignSummary{Slug: *c.CampaignSlug, Name: *c.CampaignName}
		}
		items = append(items, val)
	}

	return pagination.Result[characterDomain.Character]{Items: items, Total: total}, nil
}

// GetByPlayerInCampaigns returns the player characters linked to any of the campaigns,
// ordered by name ignoring case and accents, then id.
func (r *CharacterRepositoryPG) GetByPlayerInCampaigns(
	userID userDomain.UserID,
	campaignIDs []campaignDomain.CampaignID,
) ([]characterDomain.Character, error) {
	ids := make([]int, 0, len(campaignIDs))
	for _, id := range campaignIDs {
		ids = append(ids, id.Value())
	}

	rows, err := r.db.Query(context.Background(), `
		SELECT cs.*
		FROM character_sheet cs
		WHERE cs.player_id = $1 AND cs.campaign_id = ANY($2)
		ORDER BY unaccent(lower(cs.name)), cs.id
	`, userID.Value(), ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	record, err := pgx.CollectRows(rows, pgx.RowToStructByName[CharacterRow])
	if err != nil {
		return nil, err
	}

	characters := make([]characterDomain.Character, 0, len(record))
	for _, c := range record {
		val, err := MapRowToDomain(c)
		if err != nil {
			return nil, err
		}
		characters = append(characters, val)
	}

	return characters, nil
}

// GetAllByCampaignID returns the campaign characters that have a player (its party)
func (r *CharacterRepositoryPG) GetAllByCampaignID(campaignID campaignDomain.CampaignID) ([]characterDomain.Character, error) {
	return r.getByCampaign(campaignID, "cs.player_id IS NOT NULL")
}

// GetPremadeByCampaignID returns the campaign unclaimed premade characters, ordered by name
func (r *CharacterRepositoryPG) GetPremadeByCampaignID(campaignID campaignDomain.CampaignID) ([]characterDomain.Character, error) {
	return r.getByCampaign(campaignID, "cs.player_id IS NULL")
}

func (r *CharacterRepositoryPG) getByCampaign(campaignID campaignDomain.CampaignID, playerCondition string) ([]characterDomain.Character, error) {
	rows, err := r.db.Query(context.Background(), `
		SELECT cs.*
		FROM character_sheet cs
		WHERE cs.campaign_id = $1 AND `+playerCondition+`
		ORDER BY unaccent(lower(cs.name)), cs.id
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

func (r *CharacterRepositoryPG) Create(name characterDomain.CharacterName, playerID userDomain.UserID, system rpgDomain.System) (characterDomain.Character, error) {
	rows, err := r.db.Query(context.Background(), `
		INSERT INTO character_sheet(name, player_id, game_system)
		VALUES($1, $2, $3)
		RETURNING *
	`, name.Value(), playerID.Value(), system.Value())
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

// CreatePremade inserts a character in the campaign with no player
func (r *CharacterRepositoryPG) CreatePremade(name characterDomain.CharacterName, campaignID campaignDomain.CampaignID, system rpgDomain.System) (characterDomain.Character, error) {
	rows, err := r.db.Query(context.Background(), `
		INSERT INTO character_sheet(name, campaign_id, game_system)
		VALUES($1, $2, $3)
		RETURNING *
	`, name.Value(), campaignID.Value(), system.Value())
	if err != nil {
		return characterDomain.Character{}, err
	}
	defer rows.Close()
	record, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[CharacterRow])
	if err != nil {
		return characterDomain.Character{}, err
	}

	return MapRowToDomain(record)
}

// UpdateDetails saves the character name and its player's table name; a new name regenerates the slug
func (r *CharacterRepositoryPG) UpdateDetails(character characterDomain.Character) (characterDomain.Character, error) {
	var playerName any
	if character.PlayerName != nil {
		playerName = character.PlayerName.Value()
	}

	rows, err := r.db.Query(context.Background(), `
		UPDATE character_sheet SET name = $1, player_name = $2
		WHERE id = $3
		RETURNING *
	`, character.Name.Value(), playerName, character.Id.Value())
	if err != nil {
		return characterDomain.Character{}, err
	}
	defer rows.Close()
	record, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[CharacterRow])
	if err != nil {
		return characterDomain.Character{}, err
	}

	return MapRowToDomain(record)
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

// ClaimPremade makes userID the player of an unclaimed premade character of the campaign,
// unless userID already has a character there. Returns nil when no row changed.
func (r *CharacterRepositoryPG) ClaimPremade(characterID characterDomain.CharacterID, campaignID campaignDomain.CampaignID, userID userDomain.UserID, playerName rpgDomain.TableName) (*characterDomain.Character, error) {
	rows, err := r.db.Query(context.Background(), `
		UPDATE character_sheet cs
		SET player_id = $3, player_name = $4
		WHERE cs.id = $1
			AND cs.campaign_id = $2
			AND cs.player_id IS NULL
			AND NOT EXISTS (
				SELECT 1 FROM character_sheet o WHERE o.campaign_id = $2 AND o.player_id = $3
			)
		RETURNING cs.*
	`, characterID.Value(), campaignID.Value(), userID.Value(), playerName.Value())
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

// UnlinkFromCampaign removes the character from the campaign, with its player's table name;
// false when it was not in it
func (r *CharacterRepositoryPG) UnlinkFromCampaign(characterID characterDomain.CharacterID, campaignID campaignDomain.CampaignID) (bool, error) {
	cmdTag, err := r.db.Exec(context.Background(), `
		UPDATE character_sheet SET campaign_id = NULL, player_name = NULL
		WHERE id = $1 AND campaign_id = $2
	`, characterID.Value(), campaignID.Value())
	if err != nil {
		return false, err
	}

	return cmdTag.RowsAffected() > 0, nil
}

func (r *CharacterRepositoryPG) UpdateCampaign(campaignID campaignDomain.CampaignID, characterSlug rpgDomain.Slug, userID userDomain.UserID, playerName rpgDomain.TableName) (*characterDomain.Character, error) {
	rows, err := r.db.Query(context.Background(), `
		UPDATE character_sheet cs
		SET campaign_id = c.id, player_name = $4
		FROM campaign c
		WHERE
			cs.slug = $2
			AND cs.campaign_id IS NULL
			AND c.id = $1
			AND cs.player_id = $3
			AND cs.game_system = c.game_system
		RETURNING cs.*
	`, campaignID.Value(), characterSlug.Value(), userID.Value(), playerName.Value())
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

// RecordAccess sets the user last access to the character to now
func (r *CharacterRepositoryPG) RecordAccess(userID userDomain.UserID, id characterDomain.CharacterID) error {
	_, err := r.db.Exec(context.Background(), `
		INSERT INTO character_access (user_id, character_id, accessed_at) VALUES ($1, $2, now())
		ON CONFLICT (user_id, character_id) DO UPDATE SET accessed_at = now()
	`, userID.Value(), id.Value())
	return err
}
