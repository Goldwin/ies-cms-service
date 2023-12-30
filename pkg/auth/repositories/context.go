package repositories

type CommandContext interface {
	AccountRepository() AccountRepository
	OtpRepository() OtpRepository
}
