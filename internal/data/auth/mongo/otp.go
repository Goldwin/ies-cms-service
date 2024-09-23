package mongo

import (
	"context"

	"github.com/Goldwin/ies-pik-cms/pkg/auth/entities"
	"github.com/Goldwin/ies-pik-cms/pkg/auth/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type otpRepositoryImpl struct {
	ctx        context.Context
	db         *mongo.Database
	collection *mongo.Collection
}

// Delete implements repositories.OtpRepository.
func (o *otpRepositoryImpl) Delete(e *entities.Otp) error {
	_, err := o.collection.DeleteOne(o.ctx, bson.M{"_id": e.EmailAddress})

	return err
}

// Get implements repositories.OtpRepository.
func (o *otpRepositoryImpl) Get(email string) (*entities.Otp, error) {
	var model OTPModel

	err := o.collection.FindOne(o.ctx, bson.M{"_id": email}).Decode(&model)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return model.toEntity(), nil
}

// List implements repositories.OtpRepository.
func (o *otpRepositoryImpl) List(emails []string) ([]*entities.Otp, error) {
	var models []OTPModel

	cursor, err := o.collection.Find(o.ctx, bson.M{"_id": bson.M{"$in": emails}})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(o.ctx)
	if err = cursor.All(o.ctx, &models); err != nil {
		return nil, err
	}
	var result []*entities.Otp
	for _, model := range models {
		result = append(result, model.toEntity())
	}
	return result, nil
}

// Save implements repositories.OtpRepository.
func (o *otpRepositoryImpl) Save(e *entities.Otp) (*entities.Otp, error) {
	model := fromOtpEntity(e)
	_, err := o.collection.UpdateByID(o.ctx, e.EmailAddress, bson.M{"$set": model}, options.Update().SetUpsert(true))
	if err != nil {
		return nil, err
	}
	return e, nil
}

func NewOtpRepository(ctx context.Context, db *mongo.Database) repositories.OtpRepository {
	return &otpRepositoryImpl{
		ctx:        ctx,
		db:         db,
		collection: db.Collection(OTPCollection),
	}
}
