package entities

import "fmt"

type PersonAttendanceSummary struct {
	PersonID             string
	ScheduleID           string
	FirstEventAttendance *Attendance
	LastEventAttendance  *Attendance
}

func (e PersonAttendanceSummary) ID() string {
	return fmt.Sprintf("%s.%s", e.PersonID, e.ScheduleID)
}
