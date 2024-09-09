package queries

import (
	"fmt"

	"github.com/Goldwin/ies-pik-cms/pkg/attendance/dto"
	"github.com/Goldwin/ies-pik-cms/pkg/common/queries"
)

type GetEventFilter struct {
	ScheduleID string `json:"scheduleId" form:"scheduleId"`
	EventID    string `json:"eventId" form:"eventId"`
}

func (f GetEventFilter) Validate() error {
	if f.EventID == "" && f.ScheduleID == "" {
		return fmt.Errorf("event id or schedule id is required")
	}
	return nil
}

type GetEventResult struct {
	Data dto.EventDTO `json:"data" form:"data"`
}

type GetEvent interface {
	Execute(filter GetEventFilter) (GetEventResult, queries.QueryErrorDetail)
}
