package mongo

import (
	"context"

	"github.com/Goldwin/ies-pik-cms/pkg/attendance/entities"
	"github.com/Goldwin/ies-pik-cms/pkg/attendance/repositories"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type PersonRepositoryImpl struct {
	ctx        context.Context
	db         *mongo.Database
	collection *mongo.Collection
}

// Get implements repositories.PersonRepository.
func (p *PersonRepositoryImpl) Get(id string) (*entities.Person, error) {
	var model PersonModel
	err := p.collection.FindOne(p.ctx, bson.M{"_id": id}).Decode(&model)

	if err == mongo.ErrNoDocuments {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}
	return model.ToEntity(), nil
}

// List implements repositories.PersonRepository.
func (p *PersonRepositoryImpl) List(idList []string) ([]*entities.Person, error) {
	cursor, err := p.collection.Find(p.ctx, bson.M{"_id": bson.M{"$in": idList}})
	if err != nil {
		return nil, err
	}
	var models []PersonModel
	if err = cursor.All(p.ctx, &models); err != nil {
		return nil, err
	}
	return lo.Map(models, func(model PersonModel, _ int) *entities.Person {
		return model.ToEntity()
	}), nil
}

func NewPersonRepository(ctx context.Context, db *mongo.Database) repositories.PersonRepository {
	return &PersonRepositoryImpl{
		ctx:        ctx,
		db:         db,
		collection: db.Collection(PersonCollection),
	}
}
