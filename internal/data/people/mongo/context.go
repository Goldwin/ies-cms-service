package mongo

import (
	"context"

	"github.com/Goldwin/ies-pik-cms/pkg/people/repositories"
	"go.mongodb.org/mongo-driver/mongo"
)

type contextImpl struct {
	householdRepository repositories.HouseholdRepository
	personRepository    repositories.PersonRepository
}

// HouseholdRepository implements repositories.Context.
func (c *contextImpl) HouseholdRepository() repositories.HouseholdRepository {
	return c.householdRepository
}

// PersonRepository implements repositories.Context.
func (c *contextImpl) PersonRepository() repositories.PersonRepository {
	return c.personRepository
}

func NewContext(ctx context.Context, db *mongo.Database) repositories.CommandContext {
	return &contextImpl{
		householdRepository: NewHouseholdRepository(ctx, db),
		personRepository:    NewPersonRepository(ctx, db),
	}
}
