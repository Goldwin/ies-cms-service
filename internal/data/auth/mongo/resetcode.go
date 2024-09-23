package mongo

import (
	"context"

	"github.com/Goldwin/ies-pik-cms/pkg/auth/entities"
	"github.com/Goldwin/ies-pik-cms/pkg/auth/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type passwordResetCodeRepositoryImpl struct {
	ctx context.Context
	db  *mongo.Database

	collection *mongo.Collection
}

// Delete implements repositories.PasswordResetCodeRepository.
func (p *passwordResetCodeRepositoryImpl) Delete(e *entities.PasswordResetCode) error {
	_, err := p.collection.DeleteOne(p.ctx, bson.M{"_id": string(e.Email)})

	return err
}

// Get implements repositories.PasswordResetCodeRepository.
func (p *passwordResetCodeRepositoryImpl) Get(email string) (*entities.PasswordResetCode, error) {
	var model PasswordResetCodeModel
	err := p.collection.FindOne(p.ctx, bson.M{"_id": email}).Decode(&model)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return model.toEntity(), nil
}

// List implements repositories.PasswordResetCodeRepository.
func (p *passwordResetCodeRepositoryImpl) List(emails []string) ([]*entities.PasswordResetCode, error) {
	var models []PasswordResetCodeModel

	cursor, err := p.collection.Find(p.ctx, bson.M{"_id": bson.M{"$in": emails}})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(p.ctx)
	if err = cursor.All(p.ctx, &models); err != nil {
		return nil, err
	}
	var result []*entities.PasswordResetCode
	for _, model := range models {
		result = append(result, model.toEntity())
	}

	return result, nil
}

// Save implements repositories.PasswordResetCodeRepository.
func (p *passwordResetCodeRepositoryImpl) Save(e *entities.PasswordResetCode) (*entities.PasswordResetCode, error) {
	model := toPasswordResetCodeModel(e)

	_, err := p.collection.UpdateByID(p.ctx, e.Email, bson.M{"$set": model}, options.Update().SetUpsert(true))
	if err != nil {
		return nil, err
	}
	return e, nil
}

func NewPasswordResetCodeRepository(ctx context.Context, db *mongo.Database) repositories.PasswordResetCodeRepository {
	return &passwordResetCodeRepositoryImpl{
		ctx:        ctx,
		db:         db,
		collection: db.Collection(ResetPasswordCodeCollection),
	}
}
