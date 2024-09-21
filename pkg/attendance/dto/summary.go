package dto

import "time"

type ActivityAttendanceSummaryDTO struct {
	Total       int            `json:"total"`
	Name        string         `json:"name"`
	TotalByType map[string]int `json:"totalByType"`
}

type EventAttendanceSummaryDTO struct {
	TotalCheckedIn  int `json:"totalCheckedIn"`
	TotalCheckedOut int `json:"totalCheckedOut"`
	TotalFirstTimer int `json:"totalFirstTimer"`
	Total           int `json:"total"`

	TotalByType        map[string]int                 `json:"totalByType"`
	AcitivitiesSummary []ActivityAttendanceSummaryDTO `json:"activitiesSummary"`

	Date time.Time `json:"date"`
	ID   string    `json:"id"`
}
