package mongo

import (
	"time"

	"github.com/Goldwin/ies-pik-cms/pkg/attendance/entities"
	"github.com/samber/lo"
)

const (
	EventScheduleCollection = "event_schedules"
	EventCollection        = "events"
)
type EventScheduleModel struct {
	ID             string                       `bson:"_id"`
	Name           string                       `bson:"name"`
	TimezoneOffset int                          `bson:"timezoneOffset"`
	Type           string                       `bson:"type"`
	Activities     []EventScheduleActivityModel `bson:"activities"`
	Date           time.Time                    `bson:"date"`
	Days           []time.Weekday               `bson:"days"`
	StartDate      time.Time                    `bson:"startDate"`
	EndDate        time.Time                    `bson:"endDate"`
}

func toEventScheduleModel(e *entities.EventSchedule) EventScheduleModel {
	return EventScheduleModel{
		ID:             e.ID,
		Name:           e.Name,
		TimezoneOffset: e.TimezoneOffset,
		Type:           string(e.Type),
		Activities: lo.Map(e.Activities, func(e entities.EventScheduleActivity, _ int) EventScheduleActivityModel {
			return toEventScheduleActivityModel(&e)
		}),
		Date:      e.Date,
		Days:      e.Days,
		StartDate: e.StartDate,
		EndDate:   e.EndDate,
	}
}

func (e *EventScheduleModel) ToEventSchedule() *entities.EventSchedule {
	return &entities.EventSchedule{
		ID:             e.ID,
		Name:           e.Name,
		TimezoneOffset: e.TimezoneOffset,
		Type:           entities.EventScheduleType(e.Type),
		Activities: lo.Map(e.Activities, func(e EventScheduleActivityModel, _ int) entities.EventScheduleActivity {
			return e.ToEventScheduleActivity()
		}),
		OneTimeEventSchedule: entities.OneTimeEventSchedule{
			Date: e.Date,
		},
		WeeklyEventSchedule: entities.WeeklyEventSchedule{
			Days: e.Days,
		},
		DailyEventSchedule: entities.DailyEventSchedule{
			StartDate: e.StartDate,
			EndDate:   e.EndDate,
		},
	}
}

type EventScheduleActivityModel struct {
	ID     string `bson:"_id"`
	Name   string `bson:"name"`
	Hour   int    `bson:"hour"`
	Minute int    `bson:"minute"`
}

func toEventScheduleActivityModel(e *entities.EventScheduleActivity) EventScheduleActivityModel {
	return EventScheduleActivityModel{
		ID:     e.ID,
		Name:   e.Name,
		Hour:   e.Hour,
		Minute: e.Minute,
	}
}

func (e *EventScheduleActivityModel) ToEventScheduleActivity() entities.EventScheduleActivity {
	return entities.EventScheduleActivity{
		ID:     e.ID,
		Name:   e.Name,
		Hour:   e.Hour,
		Minute: e.Minute,
	}
}

type EventModel struct {
	ID              string               `bson:"_id"`
	ScheduleID      string               `bson:"scheduleId"`
	EventActivities []EventActivityModel `bson:"eventActivities"`
	Date            time.Time            `bson:"date"`
}

func (e *EventModel) ToEvent() *entities.Event {
	return &entities.Event{
		ID:         e.ID,
		ScheduleID: e.ScheduleID,
		EventActivities: lo.Map(e.EventActivities, func(e EventActivityModel, _ int) entities.EventActivity {
			return e.ToEventActivity()
		}),
		Date: e.Date,
	}
}

func toEventModel(e *entities.Event) EventModel {
	return EventModel{
		ID:         e.ID,
		ScheduleID: e.ScheduleID,
		EventActivities: lo.Map(e.EventActivities, func(e entities.EventActivity, _ int) EventActivityModel {
			return toEventActivityModel(&e)
		}),
		Date: e.Date,
	}
}

type EventActivityModel struct {
	ID   string    `bson:"_id"`
	Name string    `bson:"name"`
	Time time.Time `bson:"time"`
}

func (e *EventActivityModel) ToEventActivity() entities.EventActivity {
	return entities.EventActivity{
		ID:   e.ID,
		Name: e.Name,
		Time: e.Time,
	}
}

func toEventActivityModel(e *entities.EventActivity) EventActivityModel {
	return EventActivityModel{
		ID:   e.ID,
		Name: e.Name,
		Time: e.Time,
	}
}
