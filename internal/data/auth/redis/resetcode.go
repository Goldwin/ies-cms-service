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

type passwordResetCodeRepositoryImpl struct {
	ctx        context.Context
	client     redis.UniversalClient
	txPipeline redis.Pipeliner
}

// Delete implements repositories.PasswordResetCodeRepository.
func (p *passwordResetCodeRepositoryImpl) Delete(e *entities.PasswordResetCode) error {
	return p.client.Del(p.ctx, getPasswordResetCodeKey(entities.EmailAddress(e.Email))).Err()
}

// Get implements repositories.PasswordResetCodeRepository.
func (p *passwordResetCodeRepositoryImpl) Get(email string) (*entities.PasswordResetCode, error) {
	val, err := p.client.Get(p.ctx, getPasswordResetCodeKey(entities.EmailAddress(email))).Bytes()
	if err != nil && err != redis.Nil {
		return nil, err
	}
	otp := entities.PasswordResetCode{}
	if len(val) == 0 {
		return nil, nil
	}
	err = msgpack.Unmarshal(val, &otp)
	if err != nil {
		return nil, err
	}
	return &otp, nil
}

// List implements repositories.PasswordResetCodeRepository.
func (p *passwordResetCodeRepositoryImpl) List(emails []string) ([]*entities.PasswordResetCode, error) {
	var result []*entities.PasswordResetCode
	for _, email := range emails {
		password, err := p.Get(email)
		if err != nil {
			return nil, fmt.Errorf("reset code for email %s not found", email)
		}
		result = append(result, password)
	}
	return result, nil
}

// Save implements repositories.PasswordResetCodeRepository.
func (p *passwordResetCodeRepositoryImpl) Save(e *entities.PasswordResetCode) (*entities.PasswordResetCode, error) {
	bytes, err := msgpack.Marshal(e)
	if err != nil {
		return nil, err
	}
	return e, p.txPipeline.Set(p.ctx, getPasswordKey(entities.EmailAddress(e.Email)), string(bytes), e.ExpiresAt.Sub(time.Now())).Err()
}

func NewPasswordResetCodeRepository(ctx context.Context, redisClient redis.UniversalClient, txPipeline redis.Pipeliner) repositories.PasswordResetCodeRepository {
	return &passwordResetCodeRepositoryImpl{
		client:     redisClient,
		ctx:        ctx,
		txPipeline: txPipeline,
	}
}

func getPasswordResetCodeKey(email entities.EmailAddress) string {
	return fmt.Sprintf("auth:password-reset-code:email#%s", email)
}
