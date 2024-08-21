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
	ID              string
	EventID         string
	EventActivityID string

	PersonID          string
	FirstName         string
	MiddleName        string
	LastName          string
	ProfilePictureUrl string

	SecurityCode   string
	SecurityNumber int
	CheckinTime    time.Time
	CheckoutTime   time.Time

	Type AttendanceType
}

func (a *Attendance) IsValid() string {
	for _, t := range AttendanceTypes {
		if a.Type == t {
			return ""
		}
	}
	return "Invalid attendance type"
}
