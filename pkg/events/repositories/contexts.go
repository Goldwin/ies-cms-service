package repositories

type CommandContext interface {
	EventCheckInRepository() EventCheckInRepository
	ChurchEventRepository() ChurchEventRepository
	PersonRepository() PersonRepository
}
