package commands

import (
	"github.com/Goldwin/ies-pik-cms/pkg/people/repositories"
)

type CommandContext interface {
	PersonRepository() repositories.PersonRepository
	HouseholdRepository() repositories.HouseholdRepository
}
