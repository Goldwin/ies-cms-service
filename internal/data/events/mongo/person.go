package mongo

import (
	"context"

	"github.com/Goldwin/ies-pik-cms/pkg/events/entities"
	"github.com/Goldwin/ies-pik-cms/pkg/events/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Person struct {
	ID        string `bson:"_id"`
	FirstName string `bson:"firstName"`
	LastName  string `bson:"lastName"`
}

type personRepositoryImpl struct {
	ctx context.Context
	db  *mongo.Database
}

// GetByIds implements repositories.PersonRepository.
func (p *personRepositoryImpl) GetByIds(ids ...string) ([]*entities.Person, error) {
	result := make([]*entities.Person, 0)
	res, err := p.db.Collection("persons").Find(p.ctx, bson.M{"_id": bson.M{"$in": ids}})
	if err != nil {
		return nil, err
	}
	for res.Next(p.ctx) {
		var model Person
		if err := res.Decode(&model); err != nil {
			return nil, err
		}
		result = append(result, &entities.Person{
			ID:        model.ID,
			FirstName: model.FirstName,
			LastName:  model.LastName,
		})
	}
	return result, nil
}

// Get implements PersonRepository.
func (p *personRepositoryImpl) Get(id string) (*entities.Person, error) {
	var model Person
	err := p.db.Collection("persons").FindOne(p.ctx, bson.M{"_id": id}).Decode(&model)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &entities.Person{
		ID:        model.ID,
		FirstName: model.FirstName,
		LastName:  model.LastName,
	}, nil
}

func NewPersonRepository(ctx context.Context, db *mongo.Database) repositories.PersonRepository {
	return &personRepositoryImpl{
		ctx: ctx,
		db:  db,
	}
}
