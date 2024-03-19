package redis

import (
	"context"
	"fmt"
	"time"

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

// DeleteResetToken implements repositories.PasswordRepository.
func (p *PasswordRepositoryImpl) DeleteResetToken(e entities.EmailAddress) error {
	return p.client.Del(p.ctx, getPasswordResetCodeKey(e)).Err()
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

func (p *PasswordRepositoryImpl) GetResetCode(e entities.EmailAddress) (string, error) {
	val := p.client.Get(p.ctx, getPasswordResetCodeKey(e)).Val()
	return val, nil
}

func (p *PasswordRepositoryImpl) SaveResetCode(e entities.EmailAddress, token string, ttl time.Duration) error {
	err := p.client.Set(p.ctx, getPasswordResetCodeKey(e), token, ttl).Err()
	return err
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

func getPasswordResetCodeKey(email entities.EmailAddress) string {
	return fmt.Sprintf("auth:password-reset-code:email#%s", email)
}
