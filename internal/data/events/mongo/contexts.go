package mongo

import (
	"context"

	"github.com/Goldwin/ies-pik-cms/pkg/events/repositories"
	. "github.com/Goldwin/ies-pik-cms/pkg/events/repositories"
	"go.mongodb.org/mongo-driver/mongo"
)

type commandContextImpl struct {
	churchEventRepository         ChurchEventRepository
	churchEventScheduleRepository ChurchEventScheduleRepository
	eventCheckInRepository        EventCheckInRepository
	personRepository              PersonRepository
}

// ChurchEventRepository implements CommandContext.
func (c *commandContextImpl) ChurchEventRepository() ChurchEventRepository {
	return c.churchEventRepository
}

// ChurchEventScheduleRepository implements CommandContext.
func (c *commandContextImpl) ChurchEventScheduleRepository() ChurchEventScheduleRepository {
	return c.churchEventScheduleRepository
}

// EventCheckInRepository implements CommandContext.
func (c *commandContextImpl) EventCheckInRepository() EventCheckInRepository {
	return c.eventCheckInRepository
}

// PersonRepository implements CommandContext.
func (c *commandContextImpl) PersonRepository() repositories.PersonRepository {
	return c.personRepository
}

func NewCommandContext(ctx context.Context, mongo *mongo.Database) CommandContext {
	return &commandContextImpl{
		churchEventRepository:         NewChurchEventRepository(ctx, mongo),
		churchEventScheduleRepository: NewChurchEventScheduleRepository(ctx, mongo),
		eventCheckInRepository:        NewEventCheckInRepository(ctx, mongo),
		personRepository:              NewPersonRepository(ctx, mongo),
	}
}
