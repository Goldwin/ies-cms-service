package redis

import (
	"context"

	"github.com/Goldwin/ies-pik-cms/pkg/auth/commands"
	"github.com/Goldwin/ies-pik-cms/pkg/auth/repositories"
	"github.com/redis/go-redis/v9"
)

type redisAuthContext struct {
	otpRepository               repositories.OtpRepository
	accountRepository           repositories.AccountRepository
	passwordRepository          repositories.PasswordRepository
	passwordResetCodeRepository repositories.PasswordResetCodeRepository
}

// PasswordResetCodeRepository implements commands.CommandContext.
func (c *redisAuthContext) PasswordResetCodeRepository() repositories.PasswordResetCodeRepository {
	return c.passwordResetCodeRepository
}

// PasswordRepository implements repositories.CommandContext.
func (c *redisAuthContext) PasswordRepository() repositories.PasswordRepository {
	return c.passwordRepository
}

// AccountRepository implements repositories.Context.
func (c *redisAuthContext) AccountRepository() repositories.AccountRepository {
	return c.accountRepository
}

// OtpRepository implements repositories.Context.
func (c *redisAuthContext) OtpRepository() repositories.OtpRepository {
	return c.otpRepository
}

func NewContext(ctx context.Context, client redis.UniversalClient, txPipeline redis.Pipeliner) commands.CommandContext {
	return &redisAuthContext{
		otpRepository:      NewOtpRepository(ctx, client, txPipeline),
		accountRepository:  NewAccountRepository(ctx, client, txPipeline),
		passwordRepository: NewPasswordRepository(ctx, client, txPipeline),
		passwordResetCodeRepository: NewPasswordResetCodeRepository(ctx, client, txPipeline),
	}
}
