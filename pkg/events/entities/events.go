package entities

import "time"

type ChurchEvent struct {
	ID        string
	Name      string
	StartTime time.Time
}
