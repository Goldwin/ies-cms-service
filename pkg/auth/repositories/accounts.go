//go:generate mockery --output "mocks" --all --with-expecter=true
package repositories

import "github.com/Goldwin/ies-pik-cms/pkg/auth/entities"

type AccountRepository interface {
	AddAccount(entities.Account) (*entities.Account, error)
	GetAccount(entities.EmailAddress) (*entities.Account, error)
	UpdateAccount(entities.Account) (*entities.Account, error)
}
