package entities

import "time"

type AttendanceType string

const (
	Volunteer AttendanceType = "volunteer"
	Guest     AttendanceType = "guest"
	Regular   AttendanceType = "regular"
)

var (
	AttendanceTypes = []AttendanceType{Volunteer, Guest, Regular}
)

type Attendance struct {
	ID string

	Event         *Event
	EventActivity *EventActivity

	Attendee    *Person
	CheckedInBy *Person

	SecurityCode   string
	SecurityNumber int
	CheckinTime    time.Time
	CheckoutTime   time.Time

	Type      AttendanceType
	FirstTime bool
}

func (a *Attendance) IsValid() string {
	for _, t := range AttendanceTypes {
		if a.Type == t {
			return ""
		}
	}
	return "Invalid attendance type"
}
