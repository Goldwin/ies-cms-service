package mongo

import (
	"context"

	"github.com/Goldwin/ies-pik-cms/pkg/auth/entities"
	"github.com/Goldwin/ies-pik-cms/pkg/auth/repositories"
	"go.mongodb.org/mongo-driver/mongo"
)

type otpRepositoryImpl struct {
	ctx        context.Context
	db         *mongo.Database
	collection *mongo.Collection
}

// Delete implements repositories.OtpRepository.
func (o *otpRepositoryImpl) Delete(*entities.Otp) error {
	panic("unimplemented")
}

// Get implements repositories.OtpRepository.
func (o *otpRepositoryImpl) Get(string) (*entities.Otp, error) {
	panic("unimplemented")
}

// List implements repositories.OtpRepository.
func (o *otpRepositoryImpl) List([]string) ([]*entities.Otp, error) {
	panic("unimplemented")
}

// Save implements repositories.OtpRepository.
func (o *otpRepositoryImpl) Save(*entities.Otp) (*entities.Otp, error) {
	panic("unimplemented")
}

func NewOtpRepository(ctx context.Context, db *mongo.Database) repositories.OtpRepository {
	return &otpRepositoryImpl{
		ctx:        ctx,
		db:         db,
		collection: db.Collection(OTPCollection),
	}
}
