package redis

import (
	"context"
	"fmt"

	"github.com/Goldwin/ies-pik-cms/pkg/auth/entities"
	"github.com/Goldwin/ies-pik-cms/pkg/auth/repositories"
	"github.com/redis/go-redis/v9"
	"github.com/vmihailenco/msgpack/v5"
)

type accountRepositoryImpl struct {
	ctx        context.Context
	txPipeline redis.Pipeliner
	client     redis.UniversalClient
}

// Delete implements repositories.AccountRepository.
func (a *accountRepositoryImpl) Delete(account *entities.Account) error {
	_, err := a.client.Del(a.ctx, getAccountKey(account.Email)).Result()
	return err
}

// Get implements repositories.AccountRepository.
func (a *accountRepositoryImpl) Get(email string) (*entities.Account, error) {
	var account entities.Account
	bytes, err := a.client.Get(a.ctx, getAccountKey(entities.EmailAddress(email))).Bytes()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	if len(bytes) == 0 {
		return nil, nil
	}
	msgpack.Unmarshal(bytes, &account)
	return &account, nil
}

// List implements repositories.AccountRepository.
func (a *accountRepositoryImpl) List(emails []string) ([]*entities.Account, error) {
	var result []*entities.Account

	for _, email := range emails {
		account, err := a.Get(email)
		if err != nil {
			return nil, fmt.Errorf("account %s not found", email)
		}
		result = append(result, account)
	}
	return result, nil
}

// Save implements repositories.AccountRepository.
func (a *accountRepositoryImpl) Save(account *entities.Account) (*entities.Account, error) {
	bytes, err := msgpack.Marshal(account)
	if err != nil {
		return nil, err
	}
	a.txPipeline.Set(a.ctx, getAccountKey(account.Email), bytes, 0)
	return account, nil
}

// AddAccount implements repositories.AccountRepository.
func (a *accountRepositoryImpl) AddAccount(account entities.Account) (*entities.Account, error) {
	return a.UpdateAccount(account)
}

// GetAccount implements repositories.AccountRepository.
func (a *accountRepositoryImpl) GetAccount(email entities.EmailAddress) (*entities.Account, error) {
	var account entities.Account
	bytes, err := a.client.Get(a.ctx, getAccountKey(email)).Bytes()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	if len(bytes) == 0 {
		return nil, nil
	}
	msgpack.Unmarshal(bytes, &account)
	return &account, nil
}

// UpdateAccount implements repositories.AccountRepository.
func (a *accountRepositoryImpl) UpdateAccount(account entities.Account) (*entities.Account, error) {
	bytes, err := msgpack.Marshal(account)
	if err != nil {
		return nil, err
	}
	a.txPipeline.Set(a.ctx, getAccountKey(account.Email), bytes, 0)
	return &account, nil
}

func NewAccountRepository(ctx context.Context, client redis.UniversalClient, txPipeline redis.Pipeliner) repositories.AccountRepository {
	return &accountRepositoryImpl{
		ctx:        ctx,
		txPipeline: txPipeline,
		client:     client,
	}
}

func getAccountKey(email entities.EmailAddress) string {
	return fmt.Sprintf("auth:accounts:email#%s", email)
}
