package entities

import (
	"fmt"
	"time"

	"github.com/samber/lo"
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
	Activities     []EventScheduleActivity
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

/*
CreateNextEvent returns the next event entity in the schedule
Noted that all time fields of the event and activities are in UTC
*/
func (s *EventSchedule) CreateNextEvent(targetDate time.Time) (*Event, error) {
	nextEventDate, err := s.getNextDate(targetDate.UTC().Add(time.Duration(s.TimezoneOffset) * time.Hour))

	if err != nil {
		return nil, err
	}

	return &Event{
		ID:         s.ID,
		Name:       s.Name,
		ScheduleID: s.ID,
		EventActivities: lo.Map(s.Activities, func(e EventScheduleActivity, _ int) EventActivity {
			return EventActivity{
				ID:   e.ID,
				Name: e.Name,
				Time: time.Date(nextEventDate.Year(), nextEventDate.Month(), nextEventDate.Day(), e.Hour, e.Minute, 0, 0, nextEventDate.Location()).Add(time.Duration(s.TimezoneOffset) * time.Hour),
			}
		}),
		Date: nextEventDate,
	}, nil
}

func (s *EventSchedule) getNextDate(targetDate time.Time) (time.Time, error) {
	if s.IsOneTime() {
		if !s.OneTimeEventSchedule.Date.After(targetDate) {
			return targetDate, fmt.Errorf("One Time Event schedule already ended in the past at %s", s.OneTimeEventSchedule.Date.String())
		}
		return s.OneTimeEventSchedule.Date.UTC(), nil
	}

	if s.IsWeekly() {
		nextDate := targetDate.UTC().Add(time.Duration(s.TimezoneOffset) * time.Hour)

		if !lo.Contains(s.WeeklyEventSchedule.Days, nextDate.Weekday()) {
			return nextDate, fmt.Errorf("There is no weekly schedule for date %s", nextDate.String())
		}

		return nextDate.UTC().Add(time.Duration(s.TimezoneOffset) * time.Hour), nil
	}

	if s.IsDaily() {
		if s.DailyEventSchedule.EndDate.Before(targetDate) {
			return targetDate, fmt.Errorf("Daily Event Schedule already ended at %s", s.DailyEventSchedule.EndDate.String())
		}
		return s.DailyEventSchedule.StartDate.UTC(), nil
	}
	return targetDate.UTC(), nil
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
