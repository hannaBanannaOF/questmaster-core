package datasource

import (
	enum "questmaster-core/domain/enumerations"
	models "questmaster-core/internal/app/infra/models"
	"time"
)

type SessionDataSourceInterface interface {
	GetAllByPlayerIdOrDmId(UserId string) ([]models.Session, error)
	GetOneByStartDate(StartDate time.Time) (*models.Pair[models.Session, models.SessionCalendar], error)
	GetCalendar(StartDate time.Time, EndDate time.Time) ([]models.Pair[models.Session, models.SessionCalendar], error)
	GetOne(SessionId int) (*models.Session, error)
	ToggleInPlayById(SessionId int) (*models.Session, error)
	ResolveSlug(Slug string) (*int, error)
	CreateSession(SessionName string, SessionOverview *string, TrpgSystem enum.TrpgSystem, UserId string) (*models.Session, error)
}
