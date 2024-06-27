package queries

import (
	"github.com/Goldwin/ies-pik-cms/pkg/attendance/dto"
	"github.com/Goldwin/ies-pik-cms/pkg/common/queries"
)

type ListEventAttendanceQuery struct {
	ScheduleID string `json:"schedule_id"`
	EventID    string `json:"event_id"`
	Limit      int    `json:"limit"`
	LastID     string `json:"last_id"`
}

type ListEventAttendanceResult struct {
	Data []dto.EventCheckInDTO `json:"data"`
}

type ListEventAttendance interface {
	Execute(query ListEventAttendanceQuery) (ListEventAttendanceResult, queries.QueryErrorDetail)
}
