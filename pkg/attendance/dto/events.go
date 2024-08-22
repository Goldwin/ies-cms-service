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
	Date       time.Time          `json:"date"`
}

func FromEventEntities(result *entities.Event) EventDTO {
	return EventDTO{
		ID:         result.ID,
		ScheduleID: result.ScheduleID,
		Name:       result.Name,
		Activities: lo.Map(result.EventActivities,
			func(ea *entities.EventActivity, _ int) EventActivityDTO {
				return EventActivityDTO{ID: ea.ID, Name: ea.Name, Time: ea.Time}
			}),
		Date: result.Date,
	}
}
