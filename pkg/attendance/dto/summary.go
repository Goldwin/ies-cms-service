package dto

import "time"

type ActivityAttendanceSummaryDTO struct {
	Total       int
	Name        string
	TotalByType map[string]int
}

type EventAttendanceSummaryDTO struct {
	TotalCheckedIn  int
	TotalCheckedOut int
	TotalFirstTimer int
	Total           int

	TotalByType        map[string]int
	AcitivitiesSummary []ActivityAttendanceSummaryDTO

	Date time.Time
	ID   string
}
