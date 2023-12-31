package redis

import (
	"context"
	"fmt"

	"github.com/Goldwin/ies-pik-cms/pkg/auth/entities"
	"github.com/Goldwin/ies-pik-cms/pkg/auth/repositories"
	"github.com/redis/go-redis/v9"
	"github.com/vmihailenco/msgpack/v5"
)

type PasswordRepositoryImpl struct {
	ctx        context.Context
	client     redis.UniversalClient
	txPipeline redis.Pipeliner
}

// Save implements repositories.PasswordRepository.
func (p *PasswordRepositoryImpl) Save(e entities.PasswordDetail) error {
	bytes, err := msgpack.Marshal(e)
	if err != nil {
		return err
	}
	return p.txPipeline.Set(p.ctx, getPasswordKey(e.EmailAddress), string(bytes), 0).Err()
}

// Get implements repositories.PasswordRepository.
func (p *PasswordRepositoryImpl) Get(e entities.EmailAddress) (*entities.PasswordDetail, error) {
	val, err := p.client.Get(p.ctx, getPasswordKey(e)).Bytes()
	if err != nil && err != redis.Nil {
		return nil, err
	}
	otp := entities.PasswordDetail{}
	if len(val) == 0 {
		return nil, nil
	}
	err = msgpack.Unmarshal(val, &otp)
	if err != nil {
		return nil, err
	}
	return &otp, nil
}

func NewPasswordRepository(ctx context.Context, client redis.UniversalClient, txPipeline redis.Pipeliner) repositories.PasswordRepository {
	return &PasswordRepositoryImpl{
		client:     client,
		ctx:        ctx,
		txPipeline: txPipeline,
	}
}

func getPasswordKey(email entities.EmailAddress) string {
	return fmt.Sprintf("auth:password:email#%s", email)
}
