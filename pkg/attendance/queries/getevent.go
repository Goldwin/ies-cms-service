package queries

import (
	"github.com/Goldwin/ies-pik-cms/pkg/attendance/dto"
	"github.com/Goldwin/ies-pik-cms/pkg/common/queries"
)

type GetEventQuery struct {
	ScheduleID string `json:"scheduleId" form:"scheduleId"`
	EventID    string `json:"eventId" form:"eventId"`
}

type GetEventResult struct {
	Data dto.EventDTO `json:"data" form:"data"`
}

type GetEvent interface {
	Execute(query GetEventQuery) (GetEventResult, queries.QueryErrorDetail)
}
