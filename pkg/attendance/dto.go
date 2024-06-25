package attendance

import (
	"time"

	"github.com/Goldwin/ies-pik-cms/pkg/attendance/entities"
	"github.com/samber/lo"
)

type EventScheduleDTO struct {
	ID             string                     `json:"id"`
	Name           string                     `json:"name"`
	TimezoneOffset int                        `json:"timezone_offset"`
	Type           string                     `json:"type"`
	Activities     []EventScheduleActivityDTO `json:"activities"`
	Date           time.Time                  `json:"date"`
	Days           []time.Weekday             `json:"days"`
	StartDate      time.Time                  `json:"start_date"`
	EndDate        time.Time                  `json:"end_date"`
}

func fromEntities(result *entities.EventSchedule) EventScheduleDTO {
	return EventScheduleDTO{
		ID:             result.ID,
		Name:           result.Name,
		TimezoneOffset: result.TimezoneOffset,
		Type:           string(result.Type),
		Activities: lo.Map(result.Activities,
			func(ea entities.EventScheduleActivity, _ int) EventScheduleActivityDTO {
				return EventScheduleActivityDTO{ID: ea.ID, Name: ea.Name, Hour: ea.Hour, Minute: ea.Minute}
			}),
		Date:      result.Date,
		StartDate: result.StartDate,
		EndDate:   result.EndDate}
}

type EventScheduleActivityDTO struct {
	ScheduleID string
	ID         string
	Name       string
	Hour       int
	Minute     int
}

type EventDTO struct {
	ID             string
	Name           string
	TimezoneOffset int
	Type           string
	Activities     []EventActivityDTO
	Date           time.Time
	Days           []time.Weekday
	StartDate      time.Time
	EndDate        time.Time
}

type EventActivityDTO struct {
	ID   string
	Name string
	Time time.Time
}
