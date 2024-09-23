package mongo

import (
	"context"

	"github.com/Goldwin/ies-pik-cms/pkg/auth/entities"
	"github.com/Goldwin/ies-pik-cms/pkg/auth/repositories"
	"go.mongodb.org/mongo-driver/mongo"
)

type passwordResetCodeRepositoryImpl struct {
	ctx context.Context
	db  *mongo.Database

	collection *mongo.Collection
}

// Delete implements repositories.PasswordResetCodeRepository.
func (p *passwordResetCodeRepositoryImpl) Delete(*entities.PasswordResetCode) error {
	panic("unimplemented")
}

// Get implements repositories.PasswordResetCodeRepository.
func (p *passwordResetCodeRepositoryImpl) Get(string) (*entities.PasswordResetCode, error) {
	panic("unimplemented")
}

// List implements repositories.PasswordResetCodeRepository.
func (p *passwordResetCodeRepositoryImpl) List([]string) ([]*entities.PasswordResetCode, error) {
	panic("unimplemented")
}

// Save implements repositories.PasswordResetCodeRepository.
func (p *passwordResetCodeRepositoryImpl) Save(*entities.PasswordResetCode) (*entities.PasswordResetCode, error) {
	panic("unimplemented")
}

func NewPasswordResetCodeRepository(ctx context.Context, db *mongo.Database) repositories.PasswordResetCodeRepository {
	return &passwordResetCodeRepositoryImpl{
		ctx:        ctx,
		db:         db,
		collection: db.Collection(ResetPasswordCodeCollection),
	}
}
