package mongo

import (
	"context"

	"github.com/Goldwin/ies-pik-cms/pkg/auth/entities"
	"github.com/Goldwin/ies-pik-cms/pkg/auth/repositories"
	"go.mongodb.org/mongo-driver/mongo"
)

type PasswordRepositoryImpl struct {
	ctx context.Context
	db  *mongo.Database

	collection *mongo.Collection
}

// Delete implements repositories.PasswordRepository.
func (p *PasswordRepositoryImpl) Delete(*entities.PasswordDetail) error {
	panic("unimplemented")
}

// Get implements repositories.PasswordRepository.
func (p *PasswordRepositoryImpl) Get(string) (*entities.PasswordDetail, error) {
	panic("unimplemented")
}

// List implements repositories.PasswordRepository.
func (p *PasswordRepositoryImpl) List([]string) ([]*entities.PasswordDetail, error) {
	panic("unimplemented")
}

// Save implements repositories.PasswordRepository.
func (p *PasswordRepositoryImpl) Save(*entities.PasswordDetail) (*entities.PasswordDetail, error) {
	panic("unimplemented")
}

func NewPasswordRepository(ctx context.Context, db *mongo.Database) repositories.PasswordRepository {
	return &PasswordRepositoryImpl{
		ctx: ctx,
		db:  db,
	}
}
