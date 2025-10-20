package services_v1

import (
	"log"
	vo "questmaster-core/domain/vo"
	datasource "questmaster-core/internal/app/infra/datasources"
	"questmaster-core/internal/app/infra/db"
	models "questmaster-core/internal/app/infra/models"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
)

type SessionServiceV1 struct {
	SessionDs datasource.SessionDataSourceInterface
}

func (svc *SessionServiceV1) GetAllByPlayerIdOrDmId(UserId string) []vo.SessionListItem {
	data, err := svc.SessionDs.GetAllByPlayerIdOrDmId(UserId)
	if err != nil {
		log.Panicf("Unable to get sessions: %s", err)
	}

	uuid, err := uuid.Parse(UserId)
	if err != nil {
		log.Panicf("Unable to get sessions: %s", err)
	}
	return lo.Map(data, func(model models.Session, _ int) vo.SessionListItem {
		return vo.SessionListItem{
			Slug:        model.Slug,
			Description: model.SessionName,
			Dmed:        model.DmId == uuid,
			System:      model.TrpgSystem,
		}
	})
}

func (svc *SessionServiceV1) GetMyUpcoming(UserId string) *vo.CalendarItem {
	data, err := svc.SessionDs.GetOneByStartDate(time.Now())
	if err != nil {
		log.Panicf("Unable to get session: %s", err)
	}
	uuid, err := uuid.Parse(UserId)
	if err != nil {
		log.Panicf("Unable to get session: %s", err)
	}
	if data == nil {
		return nil
	}

	return &vo.CalendarItem{
		Slug:       data.First.Slug,
		Name:       data.First.SessionName,
		Dmed:       data.First.DmId == uuid,
		System:     data.First.TrpgSystem,
		Date:       data.Second.SessionDate,
		ScheduleId: data.Second.Id,
	}
}

func (svc *SessionServiceV1) GetCalendar(StartDate time.Time, EndDate time.Time, UserId string) []vo.CalendarItem {
	data, err := svc.SessionDs.GetCalendar(StartDate, EndDate)
	if err != nil {
		log.Panicf("Unable to get sessions: %s", err)
	}

	uuid, err := uuid.Parse(UserId)
	if err != nil {
		log.Panicf("Unable to get sessions: %s", err)
	}

	return lo.Map(data, func(model db.Pair[models.Session, models.SessionCalendar], _ int) vo.CalendarItem {
		return vo.CalendarItem{
			Slug:       model.First.Slug,
			Name:       model.First.SessionName,
			Dmed:       model.First.DmId == uuid,
			System:     model.First.TrpgSystem,
			Date:       model.Second.SessionDate,
			ScheduleId: model.Second.Id,
		}
	})
}
