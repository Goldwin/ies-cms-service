package local

import (
	"github.com/Goldwin/ies-pik-cms/pkg/auth/entities"
	"github.com/Goldwin/ies-pik-cms/pkg/auth/repositories"
)

var (
	passwordMap map[string]entities.PasswordDetail = make(map[string]entities.PasswordDetail)
)

type localPasswordRepository struct {
}

// Save implements repositories.PasswordRepository.
func (*localPasswordRepository) Save(e entities.PasswordDetail) error {
	passwordMap[string(e.EmailAddress)] = e
	return nil
}

// Get implements repositories.PasswordRepository.
func (*localPasswordRepository) Get(e entities.EmailAddress) (*entities.PasswordDetail, error) {
	result, ok := passwordMap[string(e)]
	if !ok {
		return nil, nil
	}
	return &result, nil
}

func NewPasswordRepository() repositories.PasswordRepository {
	return &localPasswordRepository{}
}
