package local

import (
	"github.com/Goldwin/ies-pik-cms/pkg/auth/commands"
	"github.com/Goldwin/ies-pik-cms/pkg/auth/repositories"
)

type localAuthContext struct {
	otpRepository      repositories.OtpRepository
	accountRepository  repositories.AccountRepository
	passwordRepository repositories.PasswordRepository
}

// PasswordRepository implements repositories.CommandContext.
func (c *localAuthContext) PasswordRepository() repositories.PasswordRepository {
	return c.passwordRepository
}

// AccountRepository implements repositories.Context.
func (c *localAuthContext) AccountRepository() repositories.AccountRepository {
	return c.accountRepository
}

// OtpRepository implements repositories.Context.
func (c *localAuthContext) OtpRepository() repositories.OtpRepository {
	return c.otpRepository
}

func NewContext() commands.CommandContext {
	return &localAuthContext{
		otpRepository:     NewOtpRepository(),
		accountRepository: NewAccountRepository(),
	}
}
