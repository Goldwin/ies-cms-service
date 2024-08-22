package entities

import "time"

/*
Event is the model of actual event that is currently running, or already happened in the past
*/
type Event struct {
	ID              string
	Name            string
	ScheduleID      string
	EventActivities []*EventActivity
	Date            time.Time
}
