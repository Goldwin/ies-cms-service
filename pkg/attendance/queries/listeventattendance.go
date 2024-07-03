package queries

import (
	"github.com/Goldwin/ies-pik-cms/pkg/attendance/dto"
	"github.com/Goldwin/ies-pik-cms/pkg/common/queries"
)

type ListEventAttendanceQuery struct {
	EventID         string `json:"event_id"`
	EventActivityID string `json:"event_activity_id"`
	Limit           int    `json:"limit"`
	LastID          string `json:"last_id"`

	AttendanceTypes []string `json:"attendance_type"`
	Name            string   `json:"name"`
}

type ListEventAttendanceResult struct {
	Data []dto.EventAttendanceDTO `json:"data"`
}

type ListEventAttendance interface {
	Execute(query ListEventAttendanceQuery) (ListEventAttendanceResult, queries.QueryErrorDetail)
}
