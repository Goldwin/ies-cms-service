package attendance

import "time"

/*
Event is the model of actual event that is currently running, or already happened in the past
*/
type Event struct {
	ID              string
	Schedule        EventSchedule
	EventActivities []EventActivity
	Date            time.Time
}

const (
	EventScheduleTypeDaily   EventScheduleType = "Daily"
	EventScheduleTypeWeekly  EventScheduleType = "Weekly"
	EventScheduleTypeOneTime EventScheduleType = "OneTime"
)
