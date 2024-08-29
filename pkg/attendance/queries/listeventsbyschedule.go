package queries

import (
	"fmt"
	"time"

	"github.com/Goldwin/ies-pik-cms/pkg/attendance/dto"
	"github.com/Goldwin/ies-pik-cms/pkg/common/queries"
)

type ListEventByScheduleQuery struct {
	ScheduleID string    `json:"scheduleId" form:"scheduleId"`
	StartDate  time.Time `json:"startDate" form:"startDate"`
	EndDate    time.Time `json:"endDate" form:"endDate"`
	Limit      int       `json:"limit" form:"limit"`
	LastID     string    `json:"lastId" form:"lastId"`
}

func (query *ListEventByScheduleQuery) Validate() error {
	if query.StartDate.IsZero() {
		return fmt.Errorf("start date is required")
	}

	if query.EndDate.IsZero() {
		return fmt.Errorf("end date is required")
	}

	if query.StartDate.After(query.EndDate) {
		return fmt.Errorf("start date must be before end date")
	}

	if query.EndDate.Sub(query.StartDate) > time.Hour*24*200 {
		return fmt.Errorf("date range must be less than 200 days")
	}

	if query.Limit > 200 {
		return fmt.Errorf("limit must be less than or equal to 200")
	}

	if query.Limit <= 0 {
		return fmt.Errorf("limit must be greater than 0")
	}

	return nil
}

type ListEventByScheduleResult struct {
	Data []dto.EventDTO `json:"data" form:"data"`
}

type ListEventBySchedule interface {
	Execute(query ListEventByScheduleQuery) (ListEventByScheduleResult, queries.QueryErrorDetail)
}
