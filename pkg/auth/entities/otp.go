package entities

import "time"

type Otp struct {
	EmailAddress EmailAddress
	PasswordHash []byte
	Salt         []byte
	ExpiresAt    time.Time
}

func (o Otp) IsExpired() bool {
	return time.Now().After(o.ExpiresAt)
}
