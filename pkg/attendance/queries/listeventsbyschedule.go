package queries

import (
	"time"

	"github.com/Goldwin/ies-pik-cms/pkg/attendance/dto"
	"github.com/Goldwin/ies-pik-cms/pkg/common/queries"
)

type ListEventByScheduleQuery struct {
	ScheduleID string    `json:"scheduleId" form:"scheduleId"`
	StartDate  time.Time `json:"startDate" form:"startDate"`
	EndDate    time.Time `json:"endDate" form:"endDate"`
	Limit      int       `json:"limit" form:"limit"`
	LastDate   time.Time `json:"lastDate" form:"lastDate"`
}

type ListEventByScheduleResult struct {
	Data []dto.EventDTO `json:"data" form:"data"`
}

type ListEventBySchedule interface {
	Execute(query ListEventByScheduleQuery) (ListEventByScheduleResult, queries.QueryErrorDetail)
}
