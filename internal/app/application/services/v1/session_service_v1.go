package services_v1

import (
	"log"
	enum "questmaster-core/domain/enumerations"
	vo "questmaster-core/domain/vo"
	datasource "questmaster-core/internal/app/infra/datasources"
	models "questmaster-core/internal/app/infra/models"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
)

type SessionServiceV1 struct {
	SessionDs        datasource.SessionDataSourceInterface
	CharacterSheetDs datasource.CharacterSheetDataSourceInterface
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
			InPlay:      model.InPlay,
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

	return lo.Map(data, func(model models.Pair[models.Session, models.SessionCalendar], _ int) vo.CalendarItem {
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

func (svc *SessionServiceV1) GetSessionDetails(SessionId int, UserId string) *vo.SessionDetailItem {
	data, err := svc.SessionDs.GetOne(SessionId)
	if err != nil {
		log.Panicf("Unable to get session: %s", err)
	}

	if data == nil {
		return nil
	}

	cs, err := svc.CharacterSheetDs.GetAllBySessionId(data.Id)
	if err != nil {
		log.Panicf("Unable to get session: %s", err)
	}

	uuid, err := uuid.Parse(UserId)
	if err != nil {
		log.Panicf("Unable to get sessions: %s", err)
	}

	return &vo.SessionDetailItem{
		Id:       data.Id,
		Name:     data.SessionName,
		Overview: data.Overview,
		System:   data.TrpgSystem,
		InPlay:   data.InPlay,
		Dmed:     data.DmId == uuid,
		Characters: lo.Map(cs, func(model models.CharacterSheet, _ int) vo.SessionCharacterSheetItem {
			return vo.SessionCharacterSheetItem{
				Name:      model.CharacterName,
				MaxHp:     model.MaxHp,
				CurrentHp: model.CurrentHp,
			}
		}),
	}
}

func (svc *SessionServiceV1) ToggleInPlayById(SessionId int, UserId string) *vo.SessionDetailItem {
	data, err := svc.SessionDs.ToggleInPlayById(SessionId)
	if err != nil {
		log.Panicf("Unable to update session: %s", err)
	}
	return svc.GetSessionDetails(data.Id, UserId)
}

func (svc *SessionServiceV1) ResolveSlug(Slug string) *vo.SlugResolve {
	data, err := svc.SessionDs.ResolveSlug(Slug)
	if err != nil {
		log.Panicf("Unable to resolve slug %s: %s", Slug, err)
	}

	return &vo.SlugResolve{
		CoreId: *data,
	}
}

func (svc *SessionServiceV1) CreateSession(SessionName string, SessionOverview *string, TrpgSystem enum.TrpgSystem, UserId string) *vo.Slug {
	data, err := svc.SessionDs.CreateSession(SessionName, SessionOverview, TrpgSystem, UserId)
	if err != nil {
		log.Panicf("Unable to create session: %s", err)
	}
	return &vo.Slug{
		Slug: data.Slug,
	}
}
