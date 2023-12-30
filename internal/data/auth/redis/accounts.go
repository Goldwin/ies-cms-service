package redis

import (
	"context"
	"fmt"

	"github.com/Goldwin/ies-pik-cms/pkg/auth/entities"
	"github.com/Goldwin/ies-pik-cms/pkg/auth/repositories"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/vmihailenco/msgpack/v5"
)

type accountRepositoryImpl struct {
	ctx        context.Context
	txPipeline redis.Pipeliner
	client     redis.UniversalClient
}

// AddAccount implements repositories.AccountRepository.
func (a *accountRepositoryImpl) AddAccount(account entities.Account) (*entities.Account, error) {
	account.Person.ID = uuid.New().String()
	bytes, err := msgpack.Marshal(account)
	if err != nil {
		return nil, err
	}
	a.txPipeline.Set(a.ctx, getAccountKey(account.Email), bytes, 0)
	return &account, nil
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
	account.Person.ID = uuid.New().String()
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

func getAccountKey2(account entities.Account) string {
	return fmt.Sprintf("auth:accounts:email#%s", account.Email)
}

func getAccountKey(email entities.EmailAddress) string {
	return fmt.Sprintf("auth:accounts:email#%s", email)
}
