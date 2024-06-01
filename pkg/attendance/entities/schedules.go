package attendance

import "time"

type EventScheduleType string

/*
Event Schedule is a plan for the church event that will be run at a specific time in the future

There are 3 types of schedules:

- OneTime

- Weekly

- Daily
*/
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

/*
Event Schedule Activity describe the plan for activities that will be performed in the scheduled event
*/
type EventScheduleActivity struct {
	ID     string
	Name   string
	Hour   int
	Minute int
}
