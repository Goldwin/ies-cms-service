package queries

import (
	"time"

	"github.com/Goldwin/ies-pik-cms/pkg/attendance/dto"
	"github.com/Goldwin/ies-pik-cms/pkg/common/queries"
)

type ListEventByScheduleQuery struct {
	ScheduleID string    `json:"schedule_id" form:"schedule_id"`
	Limit      int       `json:"limit" form:"limit"`
	LastDate   time.Time `json:"last_date" form:"last_date"`
}

type ListEventByScheduleResult struct {
	Data []dto.EventDTO `json:"data" form:"data"`
}

type ListEventBySchedule interface {
	Execute(query ListEventByScheduleQuery) (ListEventByScheduleResult, queries.QueryErrorDetail)
}
