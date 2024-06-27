package queries

import (
	"github.com/Goldwin/ies-pik-cms/pkg/attendance/dto"
	"github.com/Goldwin/ies-pik-cms/pkg/common/queries"
)

type ListEventByScheduleQuery struct {
	ScheduleID string `json:"schedule_id"`
	Limit      int    `json:"limit"`
	LastID     string `json:"last_id"`
}

type ListEventByScheduleResult struct {
	Data []dto.EventScheduleDTO `json:"data"`
}

type ListEventBySchedule interface {
	Execute(query ListEventScheduleQuery) (ListEventScheduleResult, queries.QueryErrorDetail)
}
