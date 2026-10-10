package campaign

import (
	"context"
	"errors"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	campaignDomain "questmaster-core/internal/campaign/domain"
	"questmaster-core/internal/shared/pagination"
	rpgDomain "questmaster-core/internal/rpg/domain"
	userDomain "questmaster-core/internal/user/domain"
)

type CampaignRepositoryPG struct {
	db *pgxpool.Pool
}

func NewCampaignRepositoryPG(db *pgxpool.Pool) *CampaignRepositoryPG {
	return &CampaignRepositoryPG{db: db}
}

// selectCampaign selects campaigns with player_count: the number of distinct players
// with a character in the campaign, regardless of the WHERE clause appended to it.
const selectCampaign = `
	SELECT c.*,
		(SELECT COUNT(DISTINCT cs.player_id) FROM character_sheet cs WHERE cs.campaign_id = c.id) AS player_count
	FROM campaign c
`

// userCampaignsWhere restricts campaigns to the ones where userID is the DM or has a character,
// narrowed by the list filters. Returns the WHERE clause and its arguments, starting at $1.
func userCampaignsWhere(userID userDomain.UserID, filters campaignDomain.CampaignListFilters) (string, []any) {
	const isDM = "c.dm_id = $1"
	const hasCharacter = "EXISTS (SELECT 1 FROM character_sheet p WHERE p.campaign_id = c.id AND p.player_id = $1)"

	var where string
	switch {
	case filters.Role == nil:
		where = "(" + isDM + " OR " + hasCharacter + ")"
	case *filters.Role == campaignDomain.RoleDM:
		where = isDM
	default:
		where = "c.dm_id <> $1 AND " + hasCharacter
	}

	args := []any{userID.Value()}
	if filters.Status != nil {
		args = append(args, filters.Status.Value())
		where += " AND c.status = $" + strconv.Itoa(len(args))
	}

	return " WHERE " + where, args
}

// ListForUser returns a page of the campaigns where userID is the DM or has a character,
// ordered by name ignoring case and accents, then id, and the number of campaigns matching the filters.
func (r *CampaignRepositoryPG) ListForUser(
	userID userDomain.UserID,
	filters campaignDomain.CampaignListFilters,
	page pagination.Page,
) (pagination.Result[campaignDomain.Campaign], error) {
	ctx := context.Background()
	where, args := userCampaignsWhere(userID, filters)

	var total int
	if err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM campaign c"+where, args...).Scan(&total); err != nil {
		return pagination.Result[campaignDomain.Campaign]{}, err
	}

	args = append(args, page.Limit, page.Offset)
	rows, err := r.db.Query(ctx, selectCampaign+where+
		" ORDER BY unaccent(lower(c.name)), c.id"+
		" LIMIT $"+strconv.Itoa(len(args)-1)+" OFFSET $"+strconv.Itoa(len(args)), args...)
	if err != nil {
		return pagination.Result[campaignDomain.Campaign]{}, err
	}
	defer rows.Close()
	record, err := pgx.CollectRows(rows, pgx.RowToStructByName[CampaignRow])
	if err != nil {
		return pagination.Result[campaignDomain.Campaign]{}, err
	}

	items := make([]campaignDomain.Campaign, 0, len(record))
	for _, c := range record {
		val, err := MapRowToDomain(c)
		if err != nil {
			return pagination.Result[campaignDomain.Campaign]{}, err
		}
		items = append(items, val)
	}

	return pagination.Result[campaignDomain.Campaign]{Items: items, Total: total}, nil
}

func (r *CampaignRepositoryPG) FindBySlug(slug rpgDomain.Slug) (*campaignDomain.Campaign, error) {
	rows, err := r.db.Query(context.Background(), selectCampaign+`
		WHERE c.slug = $1
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
	rows, err := r.db.Query(context.Background(), selectCampaign+`
		WHERE c.id = $1
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
		RETURNING *,
			(SELECT COUNT(DISTINCT cs.player_id) FROM character_sheet cs WHERE cs.campaign_id = campaign.id) AS player_count
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
