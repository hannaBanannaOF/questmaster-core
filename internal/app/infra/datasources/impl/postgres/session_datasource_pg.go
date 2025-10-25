package datasource_pg

import (
	"context"
	enum "questmaster-core/domain/enumerations"
	"questmaster-core/internal/app/infra/models"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SessionDatasourcePG struct {
	Db *pgxpool.Pool
}

func (ds *SessionDatasourcePG) GetAllByPlayerIdOrDmId(UserId string) ([]models.Session, error) {
	var data []models.Session
	rows, err := ds.Db.Query(context.Background(), "SELECT DISTINCT s.id, s.session_name, s.dm_id, s.in_play, s.trpg_system, s.slug FROM session s LEFT JOIN character_sheet c ON c.session_id = s.id WHERE s.dm_id = $1 OR c.player_id = $1", UserId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var s models.Session
		err := rows.Scan(&s.Id, &s.SessionName, &s.DmId, &s.InPlay, &s.TrpgSystem, &s.Slug)
		if err != nil {
			return nil, err
		}
		data = append(data, s)
	}
	return data, nil
}

func (ds *SessionDatasourcePG) GetOneByStartDate(StartDate time.Time) (*models.Pair[models.Session, models.SessionCalendar], error) {
	ctx := context.Background()
	row := ds.Db.QueryRow(ctx, "SELECT s.id, s.session_name, s.dm_id, s.in_play, s.trpg_system, s.slug, sc.id, sc.session_id, sc.session_date FROM session s LEFT JOIN session_calendar sc ON sc.session_id = s.id WHERE sc.session_date >= $1", StartDate)
	var s models.Session
	var t models.SessionCalendar
	err := row.Scan(&s.Id, &s.SessionName, &s.DmId, &s.InPlay, &s.TrpgSystem, &s.Slug, &t.Id, &t.SessionId, &t.SessionDate)
	if err != nil {
		return nil, err
	}
	return &models.Pair[models.Session, models.SessionCalendar]{
		First: s, Second: t,
	}, nil
}

func (ds *SessionDatasourcePG) GetCalendar(StartDate time.Time, EndDate time.Time) ([]models.Pair[models.Session, models.SessionCalendar], error) {
	ctx := context.Background()
	var data []models.Pair[models.Session, models.SessionCalendar]
	rows, err := ds.Db.Query(ctx, "SELECT s.id, s.session_name, s.dm_id, s.in_play, s.trpg_system, s.slug, sc.id, sc.session_id, sc.session_date FROM session s LEFT JOIN session_calendar sc ON sc.session_id = s.id WHERE sc.session_date >= $1 AND sc.session_date <= $2", StartDate, EndDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var s models.Session
		var t models.SessionCalendar
		err := rows.Scan(&s.Id, &s.SessionName, &s.DmId, &s.InPlay, &s.TrpgSystem, &s.Slug, &t.Id, &t.SessionId, &t.SessionDate)
		if err != nil {
			return nil, err
		}
		data = append(data, models.Pair[models.Session, models.SessionCalendar]{
			First:  s,
			Second: t,
		})
	}
	return data, nil
}

func (ds *SessionDatasourcePG) GetOne(SessionId int) (*models.Session, error) {
	ctx := context.Background()
	row := ds.Db.QueryRow(ctx, "SELECT s.id, s.session_name, s.dm_id, s.in_play, s.trpg_system, s.slug, s.overview FROM session s WHERE s.id = $1", SessionId)
	var s models.Session
	err := row.Scan(&s.Id, &s.SessionName, &s.DmId, &s.InPlay, &s.TrpgSystem, &s.Slug, &s.Overview)
	if err != nil {
		if err != pgx.ErrNoRows {
			return nil, err
		}
		return nil, nil
	}
	return &s, nil
}

func (ds *SessionDatasourcePG) ToggleInPlayById(SessionId int) (*models.Session, error) {
	ctx := context.Background()
	row := ds.Db.QueryRow(ctx, "UPDATE session SET in_play = NOT in_play WHERE id = $1 RETURNING *", SessionId)
	var s models.Session
	err := row.Scan(&s.Id, &s.SessionName, &s.DmId, &s.InPlay, &s.TrpgSystem, &s.Slug, &s.Overview)
	if err != nil {
		if err != pgx.ErrNoRows {
			return nil, err
		}
		return nil, nil
	}
	return &s, nil
}

func (ds *SessionDatasourcePG) ResolveSlug(Slug string) (*int, error) {
	ctx := context.Background()
	row := ds.Db.QueryRow(ctx, "SELECT s.id FROM session s WHERE s.slug = $1", Slug)
	var s int
	err := row.Scan(&s)
	if err != nil {
		if err != pgx.ErrNoRows {
			return nil, err
		}
		return nil, nil
	}
	return &s, nil
}

func (ds *SessionDatasourcePG) CreateSession(SessionName string, SessionOverview *string, TrpgSystem enum.TrpgSystem, UserId string) (*models.Session, error) {
	ctx := context.Background()
	row := ds.Db.QueryRow(ctx, "INSERT INTO session(session_name, dm_id, trpg_system, overview) VALUES($1, $2, $3, $4) RETURNING *", SessionName, UserId, TrpgSystem, SessionOverview)
	var s models.Session
	err := row.Scan(&s.Id, &s.SessionName, &s.DmId, &s.InPlay, &s.TrpgSystem, &s.Slug, &s.Overview)
	if err != nil {
		if err != pgx.ErrNoRows {
			return nil, err
		}
		return nil, nil
	}
	return &s, nil
}
