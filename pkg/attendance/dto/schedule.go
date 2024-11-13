package dto

import (
	"time"

	"github.com/Goldwin/ies-pik-cms/pkg/attendance/entities"
	"github.com/samber/lo"
)

type EventScheduleDTO struct {
	ID             string                     `json:"id"`
	Name           string                     `json:"name"`
	TimezoneOffset int                        `json:"timezoneOffset"`
	Type           string                     `json:"type"`
	Activities     []EventScheduleActivityDTO `json:"activities"`
	Date           time.Time                  `json:"date"`
	Days           []time.Weekday             `json:"days"`
	StartDate      time.Time                  `json:"startDate"`
	EndDate        time.Time                  `json:"endDate"`
	StartTime      string                     `json:"startTime"`
	EndTime        string                     `json:"endTime"`
}

func FromEntities(result *entities.EventSchedule) EventScheduleDTO {
	return EventScheduleDTO{
		ID:             result.ID,
		Name:           result.Name,
		TimezoneOffset: result.TimezoneOffset,
		Type:           string(result.Type),
		Activities: lo.Map(result.Activities,
			func(ea *entities.EventScheduleActivity, _ int) EventScheduleActivityDTO {
				return EventScheduleActivityDTO{
					ID:         ea.ID,
					Name:       ea.Name,
					Hour:       ea.Hour,
					Minute:     ea.Minute,
					ScheduleID: result.ID,
					Labels: lo.Map(ea.Labels, func(label *entities.ActivityLabel, _ int) ActivityLabelDTO {
						return FromActivityLabelEntity(label)
					}),
				}
			}),
		Date:      result.Date,
		StartDate: result.StartDate,
		Days:      result.Days,
		EndDate:   result.EndDate,
		StartTime: result.StartTime.String(),
		EndTime:   result.EndTime.String(),
	}
}

type EventScheduleActivityDTO struct {
	ScheduleID string             `json:"scheduleId"`
	ID         string             `json:"id"`
	Name       string             `json:"name"`
	Hour       int                `json:"hour"`
	Minute     int                `json:"minute"`
	Labels     []ActivityLabelDTO `json:"labels"`
}


