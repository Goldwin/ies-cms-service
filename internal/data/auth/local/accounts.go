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

// Delete implements repositories.AccountRepository.
func (a *accountRepositoryImpl) Delete(account *entities.Account) error {
	delete(accounts, string(account.Email))
	return nil
}

// Get implements repositories.AccountRepository.
func (a *accountRepositoryImpl) Get(email string) (*entities.Account, error) {
	result, ok := accounts[email]
	if !ok {
		return nil, fmt.Errorf("account not found")
	}
	return &result, nil
}

// List implements repositories.AccountRepository.
func (a *accountRepositoryImpl) List(emailList []string) ([]*entities.Account, error) {

	var result []*entities.Account
	for _, email := range emailList {
		account, ok := accounts[email]
		if !ok {
			return nil, fmt.Errorf("account %s not found", email)
		}
		result = append(result, &account)
	}
	return result, nil
}

// Save implements repositories.AccountRepository.
func (a *accountRepositoryImpl) Save(account *entities.Account) (*entities.Account, error) {
	accounts[string(account.Email)] = *account
	return account, nil
}

func NewAccountRepository() repositories.AccountRepository {
	return &accountRepositoryImpl{}
}
