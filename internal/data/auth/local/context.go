package local

import "github.com/Goldwin/ies-pik-cms/pkg/auth/repositories"

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

func NewContext() repositories.CommandContext {
	return &localAuthContext{
		otpRepository:     NewOtpRepository(),
		accountRepository: NewAccountRepository(),
	}
}
