package mongo

import (
	"context"

	"github.com/Goldwin/ies-pik-cms/pkg/auth/entities"
	"github.com/Goldwin/ies-pik-cms/pkg/auth/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PasswordRepositoryImpl struct {
	ctx context.Context
	db  *mongo.Database

	collection *mongo.Collection
}

// Delete implements repositories.PasswordRepository.
func (p *PasswordRepositoryImpl) Delete(e *entities.PasswordDetail) error {
	_, err := p.collection.DeleteOne(p.ctx, bson.M{"_id": e.EmailAddress})

	return err
}

// Get implements repositories.PasswordRepository.
func (p *PasswordRepositoryImpl) Get(email string) (*entities.PasswordDetail, error) {
	var model PasswordModel
	err := p.collection.FindOne(p.ctx, bson.M{"_id": email}).Decode(&model)

	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return model.toEntity(), nil
}

// List implements repositories.PasswordRepository.
func (p *PasswordRepositoryImpl) List(emails []string) ([]*entities.PasswordDetail, error) {
	var models []PasswordModel

	cursor, err := p.collection.Find(p.ctx, bson.M{"_id": bson.M{"$in": emails}})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(p.ctx)
	if err = cursor.All(p.ctx, &models); err != nil {
		return nil, err
	}
	var result []*entities.PasswordDetail
	for _, model := range models {
		result = append(result, model.toEntity())
	}
	return result, nil
}

// Save implements repositories.PasswordRepository.
func (p *PasswordRepositoryImpl) Save(e *entities.PasswordDetail) (*entities.PasswordDetail, error) {
	var model PasswordModel

	model = fromPasswordDetailEntity(e)
	_, err := p.collection.UpdateByID(p.ctx, e.EmailAddress, bson.M{"$set": model}, options.Update().SetUpsert(true))
	if err != nil {
		return nil, err
	}
	return e, nil
}

func NewPasswordRepository(ctx context.Context, db *mongo.Database) repositories.PasswordRepository {
	return &PasswordRepositoryImpl{
		ctx: ctx,
		db:  db,

		collection: db.Collection(PasswordCollection),
	}
}
