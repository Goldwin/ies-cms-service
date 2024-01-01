package entities

import "time"

type CheckInEvent struct {
	ID        string
	Person    Person
	Event     ChurchEvent
	CheckInAt time.Time
}
