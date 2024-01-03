package commands

import "github.com/Goldwin/ies-pik-cms/pkg/events/repositories"

type CommandContext interface {
	EventCheckInRepository() repositories.EventCheckInRepository
	ChurchEventRepository() repositories.ChurchEventRepository
	PersonRepository() repositories.PersonRepository
	ChurchEventSessionRepository() repositories.ChurchEventSessionRepository
}
