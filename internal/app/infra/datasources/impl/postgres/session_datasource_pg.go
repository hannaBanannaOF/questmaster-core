package datasource_pg

import (
	"context"
	"questmaster-core/internal/app/infra/db"
	"questmaster-core/internal/app/infra/models"
	"time"

	"github.com/jackc/pgx/v5"
)

type SessionDatasourcePG struct {
	Db db.Db
}

func (ds *SessionDatasourcePG) GetAllByPlayerIdOrDmId(UserId string) ([]models.Session, error) {
	ctx := context.Background()
	db, err := ds.Db.Connect()
	if err != nil {
		return nil, err
	}
	defer db.Close(ctx)
	var data []models.Session
	rows, err := db.Query(ctx, "SELECT DISTINCT s.id, s.session_name, s.dm_id, s.in_play, s.trpg_system, s.slug FROM session s LEFT JOIN character_sheet c ON c.session_id = s.id WHERE s.dm_id = $1 OR c.player_id = $1", UserId)
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

func (ds *SessionDatasourcePG) GetOneByStartDate(StartDate time.Time) (*db.Pair[models.Session, models.SessionCalendar], error) {
	ctx := context.Background()
	conn, err := ds.Db.Connect()
	if err != nil {
		return nil, err
	}
	defer conn.Close(ctx)
	row := conn.QueryRow(ctx, "SELECT s.id, s.session_name, s.dm_id, s.in_play, s.trpg_system, s.slug, sc.id, sc.session_id, sc.session_date FROM session s LEFT JOIN session_calendar sc ON sc.session_id = s.id WHERE sc.session_date >= $1", StartDate)
	var s models.Session
	var t models.SessionCalendar
	err = row.Scan(&s.Id, &s.SessionName, &s.DmId, &s.InPlay, &s.TrpgSystem, &s.Slug, &t.Id, &t.SessionId, &t.SessionDate)
	if err != nil {
		if err != pgx.ErrNoRows {
			return nil, err
		}
		return nil, nil
	}
	return &db.Pair[models.Session, models.SessionCalendar]{
		First: s, Second: t,
	}, nil
}

func (ds *SessionDatasourcePG) GetCalendar(StartDate time.Time, EndDate time.Time) ([]db.Pair[models.Session, models.SessionCalendar], error) {
	ctx := context.Background()
	conn, err := ds.Db.Connect()
	if err != nil {
		return nil, err
	}
	defer conn.Close(ctx)
	var data []db.Pair[models.Session, models.SessionCalendar]
	rows, err := conn.Query(ctx, "SELECT s.id, s.session_name, s.dm_id, s.in_play, s.trpg_system, s.slug, sc.id, sc.session_id, sc.session_date FROM session s LEFT JOIN session_calendar sc ON sc.session_id = s.id WHERE sc.session_date >= $1 AND sc.session_date <= $2", StartDate, EndDate)
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
		data = append(data, db.Pair[models.Session, models.SessionCalendar]{
			First:  s,
			Second: t,
		})
	}
	return data, nil
}
