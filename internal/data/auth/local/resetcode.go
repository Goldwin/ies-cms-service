package local

import (
	"fmt"

	"github.com/Goldwin/ies-pik-cms/pkg/auth/entities"
	"github.com/Goldwin/ies-pik-cms/pkg/auth/repositories"
)

var (
	resetTokenMap map[string]*entities.PasswordResetCode = make(map[string]*entities.PasswordResetCode)
)

type passwordResetCodeImpl struct {
}

// Delete implements repositories.PasswordResetCodeRepository.
func (p *passwordResetCodeImpl) Delete(e *entities.PasswordResetCode) error {
	delete(resetTokenMap, e.Email)
	return nil
}

// Get implements repositories.PasswordResetCodeRepository.
func (p *passwordResetCodeImpl) Get(email string) (*entities.PasswordResetCode, error) {
	return resetTokenMap[email], nil
}

// List implements repositories.PasswordResetCodeRepository.
func (p *passwordResetCodeImpl) List(emails []string) ([]*entities.PasswordResetCode, error) {
	var result []*entities.PasswordResetCode

	for _, email := range emails {
		password, err := p.Get(email)
		if err != nil {
			return nil, fmt.Errorf("reset code for email %s not found", email)
		}
		result = append(result, password)
	}

	return result, nil
}

// Save implements repositories.PasswordResetCodeRepository.
func (p *passwordResetCodeImpl) Save(e *entities.PasswordResetCode) (*entities.PasswordResetCode, error) {
	resetTokenMap[e.Email] = e
	return e, nil
}

func NewPasswordResetCodeRepository() repositories.PasswordResetCodeRepository {
	return &passwordResetCodeImpl{}
}
