package queries

import (
	"fmt"

	"github.com/Goldwin/ies-pik-cms/pkg/attendance/dto"
	"github.com/Goldwin/ies-pik-cms/pkg/common/queries"
)

type ListEventAttendanceFilter struct {
	EventID         string `json:"eventId" form:"eventId"`
	EventActivityID string `json:"eventActivityId" form:"eventActivityId"`
	Limit           int    `json:"limit" form:"limit"`
	LastID          string `json:"lastId" form:"lastId"`

	AttendanceTypes []string `json:"attendanceType" form:"attendanceType"`
	Name            string   `json:"name" form:"name"`
}

func (f ListEventAttendanceFilter) Validate() error {
	if f.EventID == "" && f.EventActivityID == "" {
		return fmt.Errorf("event id or event activity id is required")
	}

	if f.Limit > 500 {
		return fmt.Errorf("limit must be less than or equal to 200")
	}

	return nil
}

type ListEventAttendanceResult struct {
	Data []dto.EventAttendanceDTO `json:"data" form:"data"`
}

type ListEventAttendance interface {
	Execute(filter ListEventAttendanceFilter) (ListEventAttendanceResult, queries.QueryErrorDetail)
}
