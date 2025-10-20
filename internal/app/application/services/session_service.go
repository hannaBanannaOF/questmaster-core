package services

import (
	"questmaster-core/domain/vo"
	"time"
)

type SessionServiceInterface interface {
	GetAllByPlayerIdOrDmId(UserId string) []vo.SessionListItem
	GetMyUpcoming(UserId string) *vo.CalendarItem
	GetCalendar(StartDate time.Time, EndDate time.Time, UserId string) []vo.CalendarItem
}
