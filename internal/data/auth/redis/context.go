package redis

import (
	"context"

	"github.com/Goldwin/ies-pik-cms/pkg/auth/repositories"
	"github.com/redis/go-redis/v9"
)

type localAuthContext struct {
	otpRepository     repositories.OtpRepository
	accountRepository repositories.AccountRepository
}

// AccountRepository implements repositories.Context.
func (c *localAuthContext) AccountRepository() repositories.AccountRepository {
	return c.accountRepository
}

// OtpRepository implements repositories.Context.
func (c *localAuthContext) OtpRepository() repositories.OtpRepository {
	return c.otpRepository
}

func NewContext(ctx context.Context, client redis.UniversalClient, txPipeline redis.Pipeliner) repositories.CommandContext {
	return &localAuthContext{
		otpRepository:     NewOtpRepository(ctx, client, txPipeline),
		accountRepository: NewAccountRepository(ctx, client, txPipeline),
	}
}
