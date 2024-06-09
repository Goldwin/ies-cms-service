package entities

import (
	"fmt"
	"time"
)

type EventScheduleType string

const (
	EventScheduleTypeDaily   EventScheduleType = "Daily"
	EventScheduleTypeWeekly  EventScheduleType = "Weekly"
	EventScheduleTypeOneTime EventScheduleType = "OneTime"
)

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

func (s *EventSchedule) IsOneTime() bool {
	return s.Type == EventScheduleTypeOneTime
}

func (s *EventSchedule) IsWeekly() bool {
	return s.Type == EventScheduleTypeWeekly
}

func (s *EventSchedule) IsDaily() bool {
	return s.Type == EventScheduleTypeDaily
}

func (s *EventSchedule) IsValid() string {
	if s.IsOneTime() {
		return s.OneTimeEventSchedule.IsValid()
	}
	return ""
}

type OneTimeEventSchedule struct {
	Date time.Time
}

func (s *OneTimeEventSchedule) IsValid() string {
	if !s.Date.After(time.Now()) {
		return fmt.Sprintf("Date must be in the future: %s", s.Date.String())
	}
	return ""
}

type WeeklyEventSchedule struct {
	Days []time.Weekday
}

func (s *WeeklyEventSchedule) IsValid() string {
	if len(s.Days) == 0 {
		return "Days cannot be empty"
	}
	return ""
}

type DailyEventSchedule struct {
	StartDate time.Time
	EndDate   time.Time
}

func (s *DailyEventSchedule) IsValid() string {
	if !(s.StartDate.After(time.Now()) && s.EndDate.After(s.StartDate)) {
		return fmt.Sprintf("Start Date and End Date must be in the future: Start Date: %s, End Date: %s", s.StartDate.String(), s.EndDate.String())
	}
	return ""
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
