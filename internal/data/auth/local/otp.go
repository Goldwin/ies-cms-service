package local

import (
	"errors"
	"fmt"

	"github.com/Goldwin/ies-pik-cms/pkg/auth/entities"
	"github.com/Goldwin/ies-pik-cms/pkg/auth/repositories"
)

var (
	otpMap map[string]entities.Otp = make(map[string]entities.Otp)
)

type otpRepositoryImpl struct {
}

// Delete implements repositories.OtpRepository.
func (o *otpRepositoryImpl) Delete(otp *entities.Otp) error {
	delete(otpMap, string(otp.EmailAddress))
	return nil
}

// Get implements repositories.OtpRepository.
func (o *otpRepositoryImpl) Get(email string) (*entities.Otp, error) {
	result, ok := otpMap[string(email)]
	if !ok {
		return nil, errors.New("otp not found")
	}

	return &result, nil
}

// List implements repositories.OtpRepository.
func (o *otpRepositoryImpl) List(emails []string) ([]*entities.Otp, error) {
	var result []*entities.Otp
	for _, email := range emails {
		otp, ok := otpMap[string(email)]
		if !ok {
			return nil, fmt.Errorf("otp for email %s not found", email)
		}
		result = append(result, &otp)
	}
	return result, nil
}

// Save implements repositories.OtpRepository.
func (o *otpRepositoryImpl) Save(otp *entities.Otp) (*entities.Otp, error) {
	otpMap[string(otp.EmailAddress)] = *otp
	return otp, nil
}

func NewOtpRepository() repositories.OtpRepository {
	return &otpRepositoryImpl{}
}
