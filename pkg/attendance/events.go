package attendance

import "time"

type EventScheduleType string

type EventSchedule struct {
	ID             string
	Name           string
	TimezoneOffset int
	Type           EventScheduleType
	Activities     []EventScheduleActivity
}

type OneTimeEventSchedule struct {
	Date time.Time
	EventSchedule
}

type WeeklyEventSchedule struct {
	Days []time.Weekday
	EventSchedule
}

type DailyEventSchedule struct {
	StartDate time.Time
	EndDate   time.Time
	EventSchedule
}

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
