package queries

import (
	"github.com/Goldwin/ies-pik-cms/pkg/attendance/dto"
	"github.com/Goldwin/ies-pik-cms/pkg/common/queries"
)

type GetEventFilter struct {
	ScheduleID string `json:"scheduleId" form:"scheduleId"`
	EventID    string `json:"eventId" form:"eventId"`
}

type GetEventResult struct {
	Data dto.EventDTO `json:"data" form:"data"`
}

type GetEvent interface {
	Execute(filter GetEventFilter) (GetEventResult, queries.QueryErrorDetail)
}
