package local

import (
	"errors"

	"github.com/Goldwin/ies-pik-cms/pkg/auth/entities"
	"github.com/Goldwin/ies-pik-cms/pkg/auth/repositories"
)

var (
	otpMap map[string]entities.Otp = make(map[string]entities.Otp)
)

type otpRepositoryImpl struct {
}

// AddOtp implements repositories.OtpRepository.
func (o *otpRepositoryImpl) AddOtp(otp entities.Otp) error {
	otpMap[string(otp.EmailAddress)] = otp
	return nil
}

// GetOtp implements repositories.OtpRepository.
func (o *otpRepositoryImpl) GetOtp(email entities.EmailAddress) (*entities.Otp, error) {
	result, ok := otpMap[string(email)]
	if !ok {
		return nil, errors.New("otp not found")
	}

	return &result, nil
}

// RemoveOtp implements repositories.OtpRepository.
func (o *otpRepositoryImpl) RemoveOtp(otp entities.Otp) error {
	delete(otpMap, string(otp.EmailAddress))
	return nil
}

func NewOtpRepository() repositories.OtpRepository {
	return &otpRepositoryImpl{}
}
