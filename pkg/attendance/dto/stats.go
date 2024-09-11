package dto

import "time"

type EventAttendanceCountStats struct {
	AttendanceType string `json:"attendanceType"`
	Count          int    `json:"count"`
}

type EventStatsDTO struct {
	ID   string    `json:"id"`
	Date time.Time `json:"date"`

	AttendanceCount []EventAttendanceCountStats `json:"attendanceCount"`
}

type EventScheduleStatsDTO struct {
	ID string `json:"id"`

	EventStats []EventStatsDTO `json:"eventStats"`
}
