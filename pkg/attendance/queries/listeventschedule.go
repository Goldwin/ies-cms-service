package queries

import (
	"github.com/Goldwin/ies-pik-cms/pkg/attendance/dto"
	"github.com/Goldwin/ies-pik-cms/pkg/common/queries"
)

type ListEventScheduleQuery struct {
}

type ListEventScheduleResult struct {
	data []dto.EventScheduleDTO
}

type ListEventSchedule interface {
	Execute(query ListEventScheduleQuery) (ListEventScheduleResult, queries.QueryErrorDetail)
}
