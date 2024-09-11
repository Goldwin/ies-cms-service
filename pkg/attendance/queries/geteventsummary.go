package queries

import (
	"fmt"

	"github.com/Goldwin/ies-pik-cms/pkg/attendance/dto"
	"github.com/Goldwin/ies-pik-cms/pkg/common/queries"
)

type GetEventAttendanceSummaryFilter struct {
	EventID string `json:"eventId" form:"eventId"`
}

func (f GetEventAttendanceSummaryFilter) Validate() error {
	if f.EventID == "" {
		return fmt.Errorf("event id is required")
	}
	return nil
}

type GetEventAttendanceSummaryResult struct {
	Data dto.EventAttendanceSummaryDTO `json:"data" form:"data"`
}

type GetEventAttendanceSummary interface {
	Execute(filter GetEventAttendanceSummaryFilter) (GetEventAttendanceSummaryResult, queries.QueryErrorDetail)
}
