package mongo

import (
	"context"

	"github.com/Goldwin/ies-pik-cms/pkg/people/entities"
	"github.com/Goldwin/ies-pik-cms/pkg/people/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type personRepositoryImpl struct {
	ctx        context.Context
	db         *mongo.Database
	collection *mongo.Collection
}

// Delete implements repositories.PersonRepository.
func (p *personRepositoryImpl) Delete(e *entities.Person) error {
	_, err := p.collection.DeleteOne(p.ctx, bson.M{"_id": e.ID})
	return err
}

// List implements repositories.PersonRepository.
func (p *personRepositoryImpl) List(id []string) ([]*entities.Person, error) {
	models := make([]Person, 0)
	if len(id) == 0 {
		return make([]*entities.Person, 0), nil
	}
	result, err := p.collection.Find(p.ctx, bson.M{"_id": bson.M{"$in": id}})
	if err != nil {
		return nil, err
	}
	defer result.Close(p.ctx)
	err = result.All(p.ctx, &models)
	if err != nil {
		return nil, err
	}
	entities := make([]*entities.Person, len(models))
	for i, model := range models {
		entities[i] = model.toEntity()
	}
	return entities, nil
}

// Save implements repositories.PersonRepository.
func (p *personRepositoryImpl) Save(e *entities.Person) (*entities.Person, error) {
	pp := toPersonModel(e)
	_, err := p.collection.UpdateByID(p.ctx, e.ID, bson.M{"$set": pp}, options.Update().SetUpsert(true))
	if err != nil {
		return nil, err
	}

	return e, nil
}

// Get implements repositories.PersonRepository.
func (p *personRepositoryImpl) Get(id string) (*entities.Person, error) {
	var person Person
	err := p.collection.FindOne(p.ctx, bson.M{"_id": id}).Decode(&person)
	if err != nil && err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	result := person.toEntity()
	return result, nil
}

// GetByEmail implements repositories.PersonRepository.
func (p *personRepositoryImpl) GetByEmail(email entities.EmailAddress) (*entities.Person, error) {
	var person Person
	err := p.collection.FindOne(p.ctx, bson.M{"emailAddress": string(email)}).Decode(&person)
	if err != nil && err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	result := person.toEntity()
	return result, nil
}

// ListByID implements repositories.PersonRepository.
func (p *personRepositoryImpl) ListByID(id []string) ([]entities.Person, error) {
	models := make([]Person, 0)
	if len(id) == 0 {
		return make([]entities.Person, 0), nil
	}
	result, err := p.collection.Find(p.ctx, bson.M{"_id": bson.M{"$in": id}})
	if err != nil {
		return nil, err
	}
	defer result.Close(p.ctx)
	err = result.All(p.ctx, &models)
	if err != nil {
		return nil, err
	}
	entities := make([]entities.Person, len(models))
	for i, model := range models {
		entities[i] = *model.toEntity()
	}
	return entities, nil
}

func NewPersonRepository(ctx context.Context, db *mongo.Database) repositories.PersonRepository {
	return &personRepositoryImpl{
		ctx:        ctx,
		db:         db,
		collection: db.Collection(personCollectionName),
	}
}
