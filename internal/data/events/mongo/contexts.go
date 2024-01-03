package mongo

import (
	"context"

	"github.com/Goldwin/ies-pik-cms/pkg/events/commands"
	"github.com/Goldwin/ies-pik-cms/pkg/events/queries"
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

func NewCommandContext(ctx context.Context, mongo *mongo.Database) commands.CommandContext {
	return &commandContextImpl{
		churchEventRepository:  NewChurchEventRepository(ctx, mongo),
		eventCheckInRepository: NewEventCheckInRepository(ctx, mongo),
		personRepository:       NewPersonRepository(ctx, mongo),
	}
}

type queryContextImpl struct {
	searchEvent queries.SearchEvent
}

// SearchEvent implements queries.QueryContext.
func (q *queryContextImpl) SearchEvent() queries.SearchEvent {
	return q.searchEvent
}

func NewQueryContext(ctx context.Context, db *mongo.Database) queries.QueryContext {
	return &queryContextImpl{
		searchEvent: SearchEvent(ctx, db),
	}
}
