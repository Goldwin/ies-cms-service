package dto

import "time"

type CheckInInput struct {
	PersonIDs []string `json:"person_ids"`
	EventID   string   `json:"event_id"`
	CheckerID string   `json:"checker_id"`
}
type CheckInEvent struct {
	ID        string    `json:"id"`
	Person    Person    `json:"person"`
	SessionID string    `json:"session_id"`
	CheckInAt time.Time `json:"check_in_at"`
}
