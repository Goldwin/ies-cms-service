package dto

import "time"

type CheckInInput struct {
	PersonID string
	EventID  string
}
type CheckInEvent struct {
	ID        string
	Person    Person
	Event     ChurchEvent
	CheckInAt time.Time
}