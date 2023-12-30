package repositories

import "github.com/Goldwin/ies-pik-cms/pkg/auth/entities"

type OtpRepository interface {
	AddOtp(entities.Otp) error
	RemoveOtp(entities.Otp) error
	GetOtp(entities.EmailAddress) (*entities.Otp, error)
}
