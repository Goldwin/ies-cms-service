package dto

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

func FromEntities(result *entities.EventSchedule) EventScheduleDTO {
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
	ScheduleID string `json:"schedule_id"`
	ID         string `json:"id"`
	Name       string `json:"name"`
	Hour       int    `json:"hour"`
	Minute     int    `json:"minute"`
}

type EventDTO struct {
	ID             string             `json:"id"`
	ScheduleID     string             `json:"schedule_id"`
	Name           string             `json:"name"`
	Activities     []EventActivityDTO `json:"activities"`
	Date           time.Time          `json:"date"`
}

type EventActivityDTO struct {
	ID   string    `json:"id"`
	Name string    `json:"name"`
	Time time.Time `json:"time"`
}

type EventCheckInDTO struct {
	ID                string           `json:"id"`
	ScheduleID        string           `json:"schedule_id"`
	EventID           string           `json:"event_id"`
	Activity          EventActivityDTO `json:"activity"`
	PersonID          string           `json:"person_id"`
	FirstName         string           `json:"first_name"`
	MiddleName        string           `json:"middle_name"`
	LastName          string           `json:"last_name"`
	ProfilePictureURL string           `json:"profile_picture_url"`
	SecurityCode      string           `json:"security_code"`
	SecurityNumber    int              `json:"security_number"`
	CheckinTime       time.Time        `json:"checkin_time"`
	AttendanceType    string           `json:"attendance_type"`
}
