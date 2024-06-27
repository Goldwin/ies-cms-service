package queries

import (
	"github.com/Goldwin/ies-pik-cms/pkg/attendance/dto"
	"github.com/Goldwin/ies-pik-cms/pkg/common/queries"
)

type ListEventScheduleQuery struct {
	Limit  int    `json:"limit"`
	LastID string `json:"last_id"`
}

type ListEventScheduleResult struct {
	Data []dto.EventScheduleDTO `json:"data"`
}

type ListEventSchedule interface {
	Execute(query ListEventScheduleQuery) (ListEventScheduleResult, queries.QueryErrorDetail)
}
