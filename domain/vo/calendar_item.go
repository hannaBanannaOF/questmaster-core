package vo

import (
	enum "questmaster-core/domain/enumerations"
	"time"
)

type CalendarItem struct {
	Slug       string          `json:"slug"`
	Name       string          `json:"name"`
	Dmed       bool            `json:"dmed"`
	System     enum.TrpgSystem `json:"system"`
	Date       time.Time       `json:"date"`
	ScheduleId int             `json:"scheduleId"`
}
