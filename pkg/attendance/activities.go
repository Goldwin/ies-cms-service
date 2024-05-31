package attendance

import "time"

type EventScheduleActivity struct {
	ID     string
	Name   string
	Hour   int
	Minute int
}

type EventActivity struct {
	ID   string
	Name string
	Time time.Time
}

