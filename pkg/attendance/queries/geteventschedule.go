package queries

import (
	"github.com/Goldwin/ies-pik-cms/pkg/attendance/dto"
	"github.com/Goldwin/ies-pik-cms/pkg/common/queries"
)

type GetEventScheduleQuery struct {
	ScheduleID string `json:"schedule_id"`
}

type GetEventScheduleResult struct {
	Data dto.EventScheduleDTO `json:"data"`
}

type GetEventSchedule interface {
	Execute(query GetEventScheduleQuery) (GetEventScheduleResult, queries.QueryErrorDetail)
}
