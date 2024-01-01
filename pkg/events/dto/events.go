package dto

import "time"

type CreateEventInput struct {
	ID     string
	Year   int
	Month  int
	Day    int
	Hours  int
	Minute int
}

type ChurchEvent struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	StartTime time.Time `json:"start_time"`
}
