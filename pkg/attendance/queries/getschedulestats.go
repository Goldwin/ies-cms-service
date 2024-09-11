package queries

import (
	"fmt"

	"github.com/Goldwin/ies-pik-cms/pkg/attendance/dto"
	"github.com/Goldwin/ies-pik-cms/pkg/common/queries"
)

type GetEventScheduleStatsFilter struct {
	ScheduleID string `json:"scheduleId" form:"scheduleId"`
}

func (f GetEventScheduleStatsFilter) Validate() error {
	if f.ScheduleID == "" {
		return fmt.Errorf("schedule id is required")
	}
	return nil
}

type GetEventScheduleStatsResult struct {
	Data dto.EventScheduleStatsDTO `json:"data" form:"data"`
}

type GetEventScheduleStats interface {
	Execute(filter GetEventScheduleStatsFilter) (GetEventScheduleStatsResult, queries.QueryErrorDetail)
}
