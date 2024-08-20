package entities

import "time"

type AttendanceType string

const (
	Volunteer AttendanceType = "volunteer"
	Guest     AttendanceType = "guest"
	Regular   AttendanceType = "regular"
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
