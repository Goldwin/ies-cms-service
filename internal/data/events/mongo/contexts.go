package mongo

import (
	"context"

	. "github.com/Goldwin/ies-pik-cms/pkg/events/repositories"
	"go.mongodb.org/mongo-driver/mongo"
)

type commandContextImpl struct {
	churchEventRepository  ChurchEventRepository
	eventCheckInRepository EventCheckInRepository
	personRepository       PersonRepository
}

// ChurchEventRepository implements repositories.CommandContext.
func (c *commandContextImpl) ChurchEventRepository() ChurchEventRepository {
	return c.churchEventRepository
}

// ChurchEventSessionRepository implements repositories.CommandContext.
func (c *commandContextImpl) ChurchEventSessionRepository() ChurchEventSessionRepository {
	return c.ChurchEventSessionRepository()
}

// EventCheckInRepository implements repositories.CommandContext.
func (c *commandContextImpl) EventCheckInRepository() EventCheckInRepository {
	return c.eventCheckInRepository
}

// PersonRepository implements repositories.CommandContext.
func (c *commandContextImpl) PersonRepository() PersonRepository {
	return c.personRepository
}

func NewCommandContext(ctx context.Context, mongo *mongo.Database) CommandContext {
	return &commandContextImpl{
		churchEventRepository:  NewChurchEventRepository(ctx, mongo),
		eventCheckInRepository: NewEventCheckInRepository(ctx, mongo),
		personRepository:       NewPersonRepository(ctx, mongo),
	}
}
