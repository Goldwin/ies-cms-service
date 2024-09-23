package local

import (
	"fmt"
	"time"

	"github.com/Goldwin/ies-pik-cms/pkg/auth/entities"
	"github.com/Goldwin/ies-pik-cms/pkg/auth/repositories"
)

var (
	passwordMap map[string]entities.PasswordDetail = make(map[string]entities.PasswordDetail)
	resetToken  map[string]string                  = make(map[string]string)
)

type localPasswordRepository struct {
}

// Delete implements repositories.PasswordRepository.
func (l *localPasswordRepository) Delete(e *entities.PasswordDetail) error {
	delete(passwordMap, string(e.EmailAddress))
	return nil
}

// List implements repositories.PasswordRepository.
func (l *localPasswordRepository) List(emailList []string) ([]*entities.PasswordDetail, error) {
	var result []*entities.PasswordDetail
	for _, email := range emailList {
		password, ok := passwordMap[email]
		if !ok {
			return nil, fmt.Errorf("password for email %s not found", email)
		}
		result = append(result, &password)
	}

	return result, nil
}

// DeleteResetToken implements repositories.PasswordRepository.
func (l *localPasswordRepository) DeleteResetToken(e entities.EmailAddress) error {
	delete(resetToken, string(e))
	return nil
}

// GetResetCode implements repositories.PasswordRepository.
func (l *localPasswordRepository) GetResetCode(e entities.EmailAddress) (string, error) {
	return resetToken[string(e)], nil
}

// SaveResetCode implements repositories.PasswordRepository.
func (l *localPasswordRepository) SaveResetCode(e entities.EmailAddress, token string, ttl time.Duration) error {
	resetToken[string(e)] = token
	return nil
}

// Save implements repositories.PasswordRepository.
func (*localPasswordRepository) Save(e *entities.PasswordDetail) (*entities.PasswordDetail, error) {
	passwordMap[string(e.EmailAddress)] = *e
	return e, nil
}

// Get implements repositories.PasswordRepository.
func (*localPasswordRepository) Get(e string) (*entities.PasswordDetail, error) {
	result, ok := passwordMap[e]
	if !ok {
		return nil, nil
	}
	return &result, nil
}

func NewPasswordRepository() repositories.PasswordRepository {
	return &localPasswordRepository{}
}
