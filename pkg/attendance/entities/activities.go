package entities

import "time"

/*
Event Activity is the actual activities that is associated with a distinct event
*/
type EventActivity struct {
	ID   string
	Name string
	Time time.Time
}
