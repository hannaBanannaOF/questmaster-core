package services

import (
	enum "questmaster-core/domain/enumerations"
	"questmaster-core/domain/vo"
	"time"
)

type SessionServiceInterface interface {
	GetAllByPlayerIdOrDmId(UserId string) []vo.SessionListItem
	GetMyUpcoming(UserId string) *vo.CalendarItem
	GetCalendar(StartDate time.Time, EndDate time.Time, UserId string) []vo.CalendarItem
	GetSessionDetails(SessionId int, UserId string) *vo.SessionDetailItem
	ToggleInPlayById(SessionId int, UserId string) *vo.SessionDetailItem
	ResolveSlug(Slug string) *vo.SlugResolve
	CreateSession(SessionName string, SessionOverview *string, TrpgSystem enum.TrpgSystem, UserId string) *vo.Slug
}
