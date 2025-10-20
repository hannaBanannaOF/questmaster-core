package models

import "time"

type SessionCalendar struct {
	Id          int       `db:"id"`
	SessionId   int       `db:"session_id"`
	SessionDate time.Time `db:"session_date"`
}
