package dto

import (
	"time"

	"github.com/Goldwin/ies-pik-cms/pkg/attendance/entities"
	"github.com/samber/lo"
)

type EventDTO struct {
	ID         string             `json:"id"`
	ScheduleID string             `json:"scheduleId"`
	Name       string             `json:"name"`
	Activities []EventActivityDTO `json:"activities"`
	StartDate  time.Time          `json:"startDate"`
	EndDate    time.Time          `json:"endDate"`
}

func FromEventEntities(result *entities.Event) EventDTO {
	return EventDTO{
		ID:         result.ID,
		ScheduleID: result.ScheduleID,
		Name:       result.Name,
		Activities: lo.Map(result.EventActivities,
			func(ea *entities.EventActivity, _ int) EventActivityDTO {
				return EventActivityDTO{
					ID: ea.ID, 
					Name: ea.Name, 
					Time: ea.Time, 
					Labels: lo.Map(ea.Labels, func(label *entities.ActivityLabel, _ int) ActivityLabelDTO {
						return FromActivityLabelEntity(label)
					}),
				}
			}),
		StartDate: result.StartDate,
		EndDate:   result.EndDate,
	}
}

type EventActivityDTO struct {
	ID     string    `json:"id"`
	Name   string    `json:"name"`
	Time   time.Time `json:"time"`
	Labels []ActivityLabelDTO  `json:"labels"`
}
