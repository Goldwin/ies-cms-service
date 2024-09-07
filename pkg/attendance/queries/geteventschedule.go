package queries

import (
	"github.com/Goldwin/ies-pik-cms/pkg/attendance/dto"
	"github.com/Goldwin/ies-pik-cms/pkg/common/queries"
)

type GetEventScheduleFilter struct {
	ScheduleID string `json:"scheduleId" form:"scheduleId"`
}

type GetEventScheduleResult struct {
	Data dto.EventScheduleDTO `json:"data" form:"data"`
}

type GetEventSchedule interface {
	Execute(filter GetEventScheduleFilter) (GetEventScheduleResult, queries.QueryErrorDetail)
}
