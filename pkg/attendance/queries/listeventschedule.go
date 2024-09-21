package queries

import (
	"fmt"

	"github.com/Goldwin/ies-pik-cms/pkg/attendance/dto"
	"github.com/Goldwin/ies-pik-cms/pkg/common/queries"
)

type ListEventScheduleFilter struct {
	Limit  int    `json:"limit" form:"limit"`
	LastID string `json:"lastId" form:"lastId"`
}

func (query ListEventScheduleFilter) Validate() error {
	if query.Limit > 200 {
		return fmt.Errorf("limit must be less than or equal to 200")
	}

	if query.Limit <= 0 {
		return fmt.Errorf("limit must be greater than 0")
	}

	return nil
}

type ListEventScheduleResult struct {
	Data []dto.EventScheduleDTO `json:"data" form:"data"`
}

type ListEventSchedule interface {
	Execute(filter ListEventScheduleFilter) (ListEventScheduleResult, queries.QueryErrorDetail)
}
