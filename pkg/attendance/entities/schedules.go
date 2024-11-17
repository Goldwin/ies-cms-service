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
	EventScheduleTypeNone    EventScheduleType = "None"
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
	Activities     []*EventScheduleActivity
	StartTime      HourMinute
	EndTime        HourMinute
	OneTimeEventSchedule
	WeeklyEventSchedule
	DailyEventSchedule
}

// IsOneTime returns true if the schedule is a one time schedule
func (s *EventSchedule) IsOneTime() bool {
	return s.Type == EventScheduleTypeOneTime
}

// IsWeekly returns true if the schedule is a weekly schedule
func (s *EventSchedule) IsWeekly() bool {
	return s.Type == EventScheduleTypeWeekly
}

// IsDaily returns true if the schedule is a daily schedule
func (s *EventSchedule) IsDaily() bool {
	return s.Type == EventScheduleTypeDaily
}

// IsNone returns true if the schedule is none
func (s *EventSchedule) IsNone() bool {
	return s.Type == EventScheduleTypeNone
}

/*
IsValid returns an error message if the schedule is invalid
If the schedule is valid, it returns an empty string
*/
func (s *EventSchedule) IsValid() string {
	if s.IsOneTime() {
		return s.OneTimeEventSchedule.IsValid()
	}
	if s.IsWeekly() {
		return s.WeeklyEventSchedule.IsValid()
	}
	if s.IsDaily() {
		return s.DailyEventSchedule.IsValid()
	}
	if s.IsNone() {
		return ""
	}
	return "Invalid schedule type"
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
	Labels []*ActivityLabel
}

