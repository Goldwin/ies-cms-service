package repositories

import "github.com/Goldwin/ies-pik-cms/pkg/auth/entities"

type PasswordRepository interface {
	Get(entities.EmailAddress) (*entities.PasswordDetail, error)
	Save(entities.PasswordDetail) error
}
