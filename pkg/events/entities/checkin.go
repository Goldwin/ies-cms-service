package entities

import "time"

type CheckInEvent struct {
	ID        string
	Person    Person
	SessionID string
	CheckInAt time.Time
}
