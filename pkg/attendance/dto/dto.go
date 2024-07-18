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
	ScheduleID string `json:"scheduleId"`
	ID         string `json:"id"`
	Name       string `json:"name"`
	Hour       int    `json:"hour"`
	Minute     int    `json:"minute"`
}

type EventDTO struct {
	ID         string             `json:"id"`
	ScheduleID string             `json:"scheduleId"`
	Name       string             `json:"name"`
	Activities []EventActivityDTO `json:"activities"`
	Date       time.Time          `json:"date"`
}

type EventActivityDTO struct {
	ID   string    `json:"id"`
	Name string    `json:"name"`
	Time time.Time `json:"time"`
}

type EventAttendanceDTO struct {
	ID                string           `json:"id"`
	EventID           string           `json:"eventId"`
	Activity          EventActivityDTO `json:"activity"`
	PersonID          string           `json:"personId"`
	FirstName         string           `json:"firstName"`
	MiddleName        string           `json:"middleName"`
	LastName          string           `json:"lastName"`
	ProfilePictureURL string           `json:"profilePictureUrl"`
	SecurityCode      string           `json:"securityCode"`
	SecurityNumber    int              `json:"securityNumber"`
	CheckinTime       time.Time        `json:"checkinTime"`
	AttendanceType    string           `json:"attendanceType"`
}
