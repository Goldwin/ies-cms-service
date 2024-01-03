package dto

import "time"

const (
	FrequencyWeekly Frequency = "WEEKLY"
	FrequencyDaily  Frequency = "DAILY"

	MaleOnly   GenderFilter = "MALE"
	FemaleOnly GenderFilter = "FEMALE"
)

type CreateSessionInput struct {
	EventID string
}

type AgeFilter struct {
	From int
	To   int
}

type GenderFilter string

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
	ShowAt                 time.Time
	HideAt                 time.Time
}

type ChurchEventSession struct {
	ID        string
	Name      string
	SessionNo int
	StartTime time.Time
	ShowAt    time.Time
	HideAt    time.Time
}
