package local

import (
	"fmt"
	"sync/atomic"

	"github.com/Goldwin/ies-pik-cms/pkg/auth/entities"
	"github.com/Goldwin/ies-pik-cms/pkg/auth/repositories"
)

var (
	accounts map[string]entities.Account = make(map[string]entities.Account)
	id       atomic.Uint32
)

type accountRepositoryImpl struct {
}

// AddAccount implements repositories.AccountRepository.
func (a *accountRepositoryImpl) AddAccount(account entities.Account) (*entities.Account, error) {
	newId := id.Add(1)
	account.Person.ID = fmt.Sprintf("%d", newId)
	accounts[string(account.Email)] = account
	return &account, nil
}

// GetAccount implements repositories.AccountRepository.
func (a *accountRepositoryImpl) GetAccount(email entities.EmailAddress) (*entities.Account, error) {
	result, ok := accounts[string(email)]
	if !ok {
		return nil, nil
	}
	return &result, nil
}

// UpdateAccount implements repositories.AccountRepository.
func (a *accountRepositoryImpl) UpdateAccount(account entities.Account) (*entities.Account, error) {
	accounts[string(account.Email)] = account
	return &account, nil
}

func NewAccountRepository() repositories.AccountRepository {
	return &accountRepositoryImpl{}
}
