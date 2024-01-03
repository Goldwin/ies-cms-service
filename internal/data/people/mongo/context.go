package mongo

import (
	"context"

	"github.com/Goldwin/ies-pik-cms/pkg/people/commands"
	"github.com/Goldwin/ies-pik-cms/pkg/people/queries"
	"github.com/Goldwin/ies-pik-cms/pkg/people/repositories"
	"go.mongodb.org/mongo-driver/mongo"
)

type commandContextImpl struct {
	householdRepository repositories.HouseholdRepository
	personRepository    repositories.PersonRepository
}

type queryContextImpl struct {
	searchPerson queries.SearchPerson
	viewPerson   queries.ViewPerson
}

// SearchPerson implements queries.QueryContext.
func (q *queryContextImpl) SearchPerson() queries.SearchPerson {
	return q.searchPerson
}

// ViewPerson implements queries.QueryContext.
func (q *queryContextImpl) ViewPerson() queries.ViewPerson {
	return q.viewPerson
}

// HouseholdRepository implements repositories.Context.
func (c *commandContextImpl) HouseholdRepository() repositories.HouseholdRepository {
	return c.householdRepository
}

// PersonRepository implements repositories.Context.
func (c *commandContextImpl) PersonRepository() repositories.PersonRepository {
	return c.personRepository
}

func NewCommandContext(ctx context.Context, db *mongo.Database) commands.CommandContext {
	return &commandContextImpl{
		householdRepository: NewHouseholdRepository(ctx, db),
		personRepository:    NewPersonRepository(ctx, db),
	}
}

func NewQueryContext(ctx context.Context, db *mongo.Database) queries.QueryContext {
	return &queryContextImpl{
		searchPerson: SearchPerson(ctx, db),
		viewPerson:   ViewPerson(ctx, db),
	}
}
