package vo

import enum "questmaster-core/domain/enumerations"

type CreateSession struct {
	SessionName     string          `json:"sessionName"`
	SessionOverview *string         `json:"sessionOverview"`
	TrpgSystem      enum.TrpgSystem `json:"system"`
}
