package mongo

import (
	"context"

	"github.com/Goldwin/ies-pik-cms/pkg/auth/entities"
	"github.com/Goldwin/ies-pik-cms/pkg/auth/repositories"
	"go.mongodb.org/mongo-driver/mongo"
)

type accountRepositoryImpl struct {
	ctx context.Context
	db  *mongo.Database

	collection *mongo.Collection
}

// Delete implements repositories.AccountRepository.
func (a *accountRepositoryImpl) Delete(*entities.Account) error {
	panic("unimplemented")
}

// Get implements repositories.AccountRepository.
func (a *accountRepositoryImpl) Get(string) (*entities.Account, error) {
	panic("unimplemented")
}

// List implements repositories.AccountRepository.
func (a *accountRepositoryImpl) List([]string) ([]*entities.Account, error) {
	panic("unimplemented")
}

// Save implements repositories.AccountRepository.
func (a *accountRepositoryImpl) Save(*entities.Account) (*entities.Account, error) {
	panic("unimplemented")
}

func NewAccountRepository(ctx context.Context, db *mongo.Database) repositories.AccountRepository {
	return &accountRepositoryImpl{
		ctx:        ctx,
		db:         db,
		collection: db.Collection(AccountCollection),
	}
}
