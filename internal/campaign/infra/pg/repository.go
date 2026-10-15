package campaign

import (
	"context"
	"errors"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	campaignDomain "questmaster-core/internal/campaign/domain"
	rpgDomain "questmaster-core/internal/rpg/domain"
	"questmaster-core/internal/shared/pagination"
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
const playerCountColumn = "(SELECT COUNT(DISTINCT cs.player_id) FROM character_sheet cs WHERE cs.campaign_id = c.id) AS player_count"

const selectCampaign = "SELECT c.*, " + playerCountColumn + " FROM campaign c"

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
	if filters.Search != nil {
		args = append(args, *filters.Search)
		where += " AND " + pagination.NameContains("c.name", len(args))
	}
	if filters.WithoutPlayers {
		where += " AND NOT EXISTS (SELECT 1 FROM character_sheet wp WHERE wp.campaign_id = c.id AND wp.player_id IS NOT NULL)"
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
	// $1 is the requester in userCampaignsWhere too
	rows, err := r.db.Query(ctx, "SELECT c.*, "+playerCountColumn+", a.accessed_at AS last_accessed_at"+
		" FROM campaign c LEFT JOIN campaign_access a ON a.user_id = $1 AND a.campaign_id = c.id"+where+
		page.Sort.OrderBy("a.accessed_at", "c.name", "c.id")+
		" LIMIT $"+strconv.Itoa(len(args)-1)+" OFFSET $"+strconv.Itoa(len(args)), args...)
	if err != nil {
		return pagination.Result[campaignDomain.Campaign]{}, err
	}
	defer rows.Close()
	record, err := pgx.CollectRows(rows, pgx.RowToStructByName[CampaignListRow])
	if err != nil {
		return pagination.Result[campaignDomain.Campaign]{}, err
	}

	items := make([]campaignDomain.Campaign, 0, len(record))
	for _, c := range record {
		val, err := MapRowToDomain(c.CampaignRow)
		if err != nil {
			return pagination.Result[campaignDomain.Campaign]{}, err
		}
		val.LastAccessedAt = c.LastAccessedAt
		items = append(items, val)
	}

	return pagination.Result[campaignDomain.Campaign]{Items: items, Total: total}, nil
}

// CountByStatusForUser counts the campaigns where userID is the DM or has a character, by status.
// Statuses without campaigns are absent from the result.
func (r *CampaignRepositoryPG) CountByStatusForUser(
	userID userDomain.UserID,
	filters campaignDomain.CampaignListFilters,
) (map[campaignDomain.CampaignStatus]int, error) {
	where, args := userCampaignsWhere(userID, filters)
	rows, err := r.db.Query(context.Background(), "SELECT c.status, COUNT(*) FROM campaign c"+where+" GROUP BY c.status", args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	counts := make(map[campaignDomain.CampaignStatus]int)
	for rows.Next() {
		var rawStatus string
		var count int
		if err := rows.Scan(&rawStatus, &count); err != nil {
			return nil, err
		}
		status, err := campaignDomain.NewCampaignStatus(rawStatus)
		if err != nil {
			return nil, err
		}
		counts[status] = count
	}

	return counts, rows.Err()
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

func (r *CampaignRepositoryPG) Create(Name campaignDomain.CampaignName, Overview *campaignDomain.CampaignOverview, DmID userDomain.UserID, System rpgDomain.System, DmName rpgDomain.TableName) (campaignDomain.Campaign, error) {
	var overview any
	if Overview != nil {
		overview = Overview.Value()
	} else {
		overview = nil
	}
	rows, err := r.db.Query(context.Background(), `
		INSERT INTO campaign(name, dm_id, game_system, overview, dm_name)
		VALUES($1, $2, $3, $4, $5) RETURNING *, 0 as player_count
	`, Name.Value(), DmID.Value(), System.Value(), overview, DmName.Value())
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

// DeleteById deletes the campaign and its unclaimed premade characters, and clears the
// players' table names, in one transaction.
// The foreign keys (migration 0005) unlink the other characters and remove the invite.
func (r *CampaignRepositoryPG) DeleteById(id campaignDomain.CampaignID) (bool, error) {
	ctx := context.Background()
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return false, err
	}
	defer tx.Rollback(ctx) // no-op after Commit

	if _, err := tx.Exec(ctx, `
		DELETE FROM character_sheet
		WHERE campaign_id = $1 AND player_id IS NULL
	`, id.Value()); err != nil {
		return false, err
	}

	// The table names belong to this campaign; the foreign key then unlinks the characters
	if _, err := tx.Exec(ctx, `
		UPDATE character_sheet SET player_name = NULL
		WHERE campaign_id = $1
	`, id.Value()); err != nil {
		return false, err
	}

	cmdTag, err := tx.Exec(ctx, `
		DELETE FROM campaign
		WHERE id = $1
	`, id.Value())
	if err != nil {
		return false, err
	}

	if err := tx.Commit(ctx); err != nil {
		return false, err
	}

	return cmdTag.RowsAffected() > 0, nil
}

// RecordAccess sets the user last access to the campaign to now
func (r *CampaignRepositoryPG) RecordAccess(userID userDomain.UserID, id campaignDomain.CampaignID) error {
	_, err := r.db.Exec(context.Background(), `
		INSERT INTO campaign_access (user_id, campaign_id, accessed_at) VALUES ($1, $2, now())
		ON CONFLICT (user_id, campaign_id) DO UPDATE SET accessed_at = now()
	`, userID.Value(), id.Value())
	return err
}

func (r *CampaignRepositoryPG) UpdateDmName(id campaignDomain.CampaignID, dmName rpgDomain.TableName) error {
	_, err := r.db.Exec(context.Background(), `
		UPDATE campaign SET dm_name = $1 WHERE id = $2
	`, dmName.Value(), id.Value())
	return err
}
