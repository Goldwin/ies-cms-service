package repositories

type CommandContext interface {
	PersonRepository() PersonRepository
	HouseholdRepository() HouseholdRepository
}
