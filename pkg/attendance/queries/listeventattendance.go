package queries

import (
	"github.com/Goldwin/ies-pik-cms/pkg/attendance/dto"
	"github.com/Goldwin/ies-pik-cms/pkg/common/queries"
)

type ListEventAttendanceQuery struct {
	EventID         string `json:"eventId" form:"eventId"`
	EventActivityID string `json:"eventActivityId" form:"eventActivityId"`
	Limit           int    `json:"limit" form:"limit"`
	LastID          string `json:"lastId" form:"lastId"`

	AttendanceTypes []string `json:"attendanceType" form:"attendanceType"`
	Name            string   `json:"name" form:"name"`
}

type ListEventAttendanceResult struct {
	Data []dto.EventAttendanceDTO `json:"data" form:"data"`
}

type ListEventAttendance interface {
	Execute(query ListEventAttendanceQuery) (ListEventAttendanceResult, queries.QueryErrorDetail)
}
