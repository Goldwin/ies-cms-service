package mongo

import (
	"context"

	"github.com/Goldwin/ies-pik-cms/pkg/auth/commands"
	"github.com/Goldwin/ies-pik-cms/pkg/auth/repositories"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoAuthContext struct {
	ctx context.Context
	db  *mongo.Database

	accountRepository           repositories.AccountRepository
	otpRepository               repositories.OtpRepository
	passwordRepository          repositories.PasswordRepository
	passwordResetCodeRepository repositories.PasswordResetCodeRepository
}

// AccountRepository implements commands.CommandContext.
func (m *mongoAuthContext) AccountRepository() repositories.AccountRepository {
	if m.accountRepository == nil {
		m.accountRepository = NewAccountRepository(m.ctx, m.db)
	}
	return m.accountRepository
}

// OtpRepository implements commands.CommandContext.
func (m *mongoAuthContext) OtpRepository() repositories.OtpRepository {
	if m.otpRepository == nil {
		m.otpRepository = NewOtpRepository(m.ctx, m.db)
	}
	return m.otpRepository
}

// PasswordRepository implements commands.CommandContext.
func (m *mongoAuthContext) PasswordRepository() repositories.PasswordRepository {
	if m.passwordRepository == nil {
		m.passwordRepository = NewPasswordRepository(m.ctx, m.db)
	}
	return m.passwordRepository
}

// PasswordResetCodeRepository implements commands.CommandContext.
func (m *mongoAuthContext) PasswordResetCodeRepository() repositories.PasswordResetCodeRepository {
	if m.passwordResetCodeRepository == nil {
		m.passwordResetCodeRepository = NewPasswordResetCodeRepository(m.ctx, m.db)
	}
	return m.passwordResetCodeRepository
}

func NewCommandContext(ctx context.Context, mongoDB *mongo.Database) commands.CommandContext {
	return &mongoAuthContext{
		ctx: ctx,
		db:  mongoDB,
	}
}
