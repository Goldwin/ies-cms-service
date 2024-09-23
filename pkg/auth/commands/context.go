//go:generate mockery --output "mocks" --all --with-expecter=true
package commands

import "github.com/Goldwin/ies-pik-cms/pkg/auth/repositories"

type CommandContext interface {
	AccountRepository() repositories.AccountRepository
	OtpRepository() repositories.OtpRepository
	PasswordRepository() repositories.PasswordRepository
	PasswordResetCodeRepository() repositories.PasswordResetCodeRepository
}
