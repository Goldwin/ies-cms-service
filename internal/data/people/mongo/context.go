package mongo

import (
	"context"

	"github.com/Goldwin/ies-pik-cms/pkg/people/commands"
	"github.com/Goldwin/ies-pik-cms/pkg/people/queries"
	"github.com/Goldwin/ies-pik-cms/pkg/people/repositories"
	"go.mongodb.org/mongo-driver/mongo"
)

type commandContextImpl struct {
	ctx                 context.Context
	db                  *mongo.Database
	householdRepository repositories.HouseholdRepository
	personRepository    repositories.PersonRepository
}

type queryContextImpl struct {
	ctx                   context.Context
	db                    *mongo.Database
	searchPerson          queries.SearchPerson
	viewPerson            queries.ViewPerson
	viewHouseholdByPerson queries.ViewHouseholdByPerson
	viewPersonByEmail     queries.ViewPersonByEmail
	searchHousehold       queries.SearchHousehold
}

// SearchHousehold implements queries.QueryContext.
func (q *queryContextImpl) SearchHousehold() queries.SearchHousehold {
	return NewSearchHousehold(q.ctx, q.db)
}

// ViewPersonByEmail implements queries.QueryContext.
func (q *queryContextImpl) ViewPersonByEmail() queries.ViewPersonByEmail {
	if q.viewPersonByEmail == nil {
		q.viewPersonByEmail = ViewPersonByEmail(q.ctx, q.db)
	}
	return q.viewPersonByEmail
}

// ViewHouseholdByPerson implements queries.QueryContext.
func (q *queryContextImpl) ViewHouseholdByPerson() queries.ViewHouseholdByPerson {
	if q.viewHouseholdByPerson == nil {
		q.viewHouseholdByPerson = ViewHouseholdByPerson(q.ctx, q.db)
	}
	return q.viewHouseholdByPerson
}

// SearchPerson implements queries.QueryContext.
func (q *queryContextImpl) SearchPerson() queries.SearchPerson {
	if q.searchPerson == nil {
		q.searchPerson = SearchPerson(q.ctx, q.db)
	}
	return q.searchPerson
}

// ViewPerson implements queries.QueryContext.
func (q *queryContextImpl) ViewPerson() queries.ViewPerson {
	if q.viewPerson == nil {
		q.viewPerson = ViewPerson(q.ctx, q.db)
	}
	return q.viewPerson
}

// HouseholdRepository implements repositories.Context.
func (c *commandContextImpl) HouseholdRepository() repositories.HouseholdRepository {
	if c.householdRepository == nil {
		c.householdRepository = NewHouseholdRepository(c.ctx, c.db)
	}
	return c.householdRepository
}

// PersonRepository implements repositories.Context.
func (c *commandContextImpl) PersonRepository() repositories.PersonRepository {
	if c.personRepository == nil {
		c.personRepository = NewPersonRepository(c.ctx, c.db)
	}
	return c.personRepository
}

func NewCommandContext(ctx context.Context, db *mongo.Database) commands.CommandContext {
	return &commandContextImpl{
		ctx:                 ctx,
		db:                  db,
	}
}

func NewQueryContext(ctx context.Context, db *mongo.Database) queries.QueryContext {
	return &queryContextImpl{
		ctx:                   ctx,
		db:                    db,
	}
}
