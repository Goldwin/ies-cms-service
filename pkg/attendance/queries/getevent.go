package queries

import (
	"github.com/Goldwin/ies-pik-cms/pkg/attendance/dto"
	"github.com/Goldwin/ies-pik-cms/pkg/common/queries"
)

type GetEventQuery struct {
	ScheduleID string `json:"schedule_id"`
	EventID    string `json:"event_id"`
}

type GetEventResult struct {
	Data dto.EventDTO `json:"data"`
}

type GetEvent interface {
	Execute(query GetEventQuery) (GetEventResult, queries.QueryErrorDetail)
}
