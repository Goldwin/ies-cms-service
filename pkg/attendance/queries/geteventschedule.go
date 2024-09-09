package queries

import (
	"fmt"

	"github.com/Goldwin/ies-pik-cms/pkg/attendance/dto"
	"github.com/Goldwin/ies-pik-cms/pkg/common/queries"
)

type GetEventScheduleFilter struct {
	ScheduleID string `json:"scheduleId" form:"scheduleId"`
}

func (f GetEventScheduleFilter) Validate() error {
	if f.ScheduleID == "" {
		return fmt.Errorf("scheduleId is required")
	}
	return nil
}

type GetEventScheduleResult struct {
	Data dto.EventScheduleDTO `json:"data" form:"data"`
}

type GetEventSchedule interface {
	Execute(filter GetEventScheduleFilter) (GetEventScheduleResult, queries.QueryErrorDetail)
}
