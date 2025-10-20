package datasource

import (
	"questmaster-core/internal/app/infra/db"
	models "questmaster-core/internal/app/infra/models"
	"time"
)

type SessionDataSourceInterface interface {
	GetAllByPlayerIdOrDmId(UserId string) ([]models.Session, error)
	GetOneByStartDate(StartDate time.Time) (*db.Pair[models.Session, models.SessionCalendar], error)
	GetCalendar(StartDate time.Time, EndDate time.Time) ([]db.Pair[models.Session, models.SessionCalendar], error)
}
