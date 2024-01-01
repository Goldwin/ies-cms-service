package repositories

type CommandContext interface {
	EventCheckInRepository() EventCheckInRepository
	ChurchEventRepository() ChurchEventRepository
	ChurchEventScheduleRepository() ChurchEventScheduleRepository
	PersonRepository() PersonRepository
}
