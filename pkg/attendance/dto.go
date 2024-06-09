package attendance

import "time"

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

type EventScheduleActivityDTO struct {
	ID     string
	Name   string
	Hour   int
	Minute int
}
