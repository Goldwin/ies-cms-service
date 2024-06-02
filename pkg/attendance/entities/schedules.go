package entities

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
	OneTimeEventSchedule
	WeeklyEventSchedule
	DailyEventSchedule
}

type OneTimeEventSchedule struct {
	Date time.Time
}

type WeeklyEventSchedule struct {
	Days []time.Weekday
}

type DailyEventSchedule struct {
	StartDate time.Time
	EndDate   time.Time
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
