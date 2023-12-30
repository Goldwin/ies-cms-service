//go:generate mockery --output "mocks" --all --with-expecter=true
package repositories

type CommandContext interface {
	AccountRepository() AccountRepository
	OtpRepository() OtpRepository
	PasswordRepository() PasswordRepository
}
