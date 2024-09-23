package entities

import "time"

type PasswordResetCode struct {
	Email     string
	Code      string
	ExpiresAt time.Time
}
