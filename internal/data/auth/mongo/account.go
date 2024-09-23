package mongo

import (
	"context"

	"github.com/Goldwin/ies-pik-cms/pkg/auth/entities"
	"github.com/Goldwin/ies-pik-cms/pkg/auth/repositories"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type accountRepositoryImpl struct {
	ctx context.Context
	db  *mongo.Database

	collection *mongo.Collection
}

// Delete implements repositories.AccountRepository.
func (a *accountRepositoryImpl) Delete(e *entities.Account) error {
	_, err := a.collection.DeleteOne(a.ctx, bson.M{"_id": e.Email})
	return err
}

// Get implements repositories.AccountRepository.
func (a *accountRepositoryImpl) Get(email string) (*entities.Account, error) {
	var model AccountModel
	err := a.collection.FindOne(a.ctx, bson.M{"_id": email}).Decode(&model)

	if err == mongo.ErrNoDocuments {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}
	return model.toEntity(), nil
}

// List implements repositories.AccountRepository.
func (a *accountRepositoryImpl) List(emails []string) ([]*entities.Account, error) {
	var models []AccountModel
	cursor, err := a.collection.Find(a.ctx, bson.M{"_id": bson.M{"$in": emails}})

	if err != nil {
		return nil, err
	}

	defer cursor.Close(a.ctx)

	if err = cursor.All(a.ctx, &models); err != nil {
		return nil, err
	}
	return lo.Map(models, func(model AccountModel, _ int) *entities.Account {
		return model.toEntity()
	}), nil
}

// Save implements repositories.AccountRepository.
func (a *accountRepositoryImpl) Save(e *entities.Account) (*entities.Account, error) {
	model := fromAccountEntity(e)

	_, err := a.collection.UpdateByID(a.ctx, e.Email, bson.M{"$set": model}, options.Update().SetUpsert(true))

	if err != nil {
		return nil, err
	}
	return e, nil
}

func NewAccountRepository(ctx context.Context, db *mongo.Database) repositories.AccountRepository {
	return &accountRepositoryImpl{
		ctx:        ctx,
		db:         db,
		collection: db.Collection(AccountCollection),
	}
}
