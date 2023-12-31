package entities

import "time"

const (
	FrequencyWeekly Frequency = "WEEKLY"
	FrequencyDaily  Frequency = "DAILY"
)

type AgeFilter struct {
	From int
	To   int
}

type Location struct {
	Name      string
	AgeFilter AgeFilter
}

type Frequency string

type ChurchEvent struct {
	ID                     string
	Name                   string
	Locations              []Location
	EventFrequency         Frequency
	LatestSessionStartTime time.Time
	LatestShowAt           time.Time
	LatestHideAt           time.Time
	LatestSessionNo        int
}

type ChurchEventSession struct {
	ID        string
	Name      string
	SessionNo int
	StartTime time.Time
	ShowAt    time.Time
	HideAt    time.Time
}
