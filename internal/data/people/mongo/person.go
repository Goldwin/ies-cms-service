package mongo

import (
	"context"

	"github.com/Goldwin/ies-pik-cms/pkg/people/entities"
	"github.com/Goldwin/ies-pik-cms/pkg/people/repositories"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type personRepositoryImpl struct {
	ctx        context.Context
	db         *mongo.Database
	collection *mongo.Collection
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
	result := toPersonEntities(person)
	return &result, nil
}

// AddPerson implements repositories.PersonRepository.
func (p *personRepositoryImpl) AddPerson(person entities.Person) (*entities.Person, error) {
	if person.ID == "" {
		person.ID = uuid.NewString()
	}
	pp := toPersonMongoModel(person)
	_, err := p.collection.InsertOne(p.ctx, pp)
	if err != nil {
		return nil, err
	}

	return &person, nil
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
	result := toPersonEntities(person)
	return &result, nil
}

// ListByID implements repositories.PersonRepository.
func (p *personRepositoryImpl) ListByID(id []string) ([]entities.Person, error) {
	models := make([]Person, 0)
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
		entities[i] = toPersonEntities(model)
	}
	return entities, nil
}

// UpdatePerson implements repositories.PersonRepository.
func (p *personRepositoryImpl) UpdatePerson(person entities.Person) (*entities.Person, error) {
	model := toPersonMongoModel(person)
	_, err := p.collection.UpdateOne(p.ctx, bson.M{"_id": person.ID}, bson.M{"$set": model})
	if err != nil {
		return nil, err
	}
	return &person, nil
}

func NewPersonRepository(ctx context.Context, db *mongo.Database) repositories.PersonRepository {
	return &personRepositoryImpl{
		ctx:        ctx,
		db:         db,
		collection: db.Collection("person"),
	}
}
