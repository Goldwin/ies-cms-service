package entities

import "time"

type PasswordResetCode struct {
	Email    string
	Code     string
	ExpiryAt time.Time
}
