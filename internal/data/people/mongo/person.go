package mongo

import (
	"context"
	"fmt"

	"github.com/Goldwin/ies-pik-cms/pkg/people/entities"
	"github.com/Goldwin/ies-pik-cms/pkg/people/repositories"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Person struct {
	ID                string  `bson:"_id"`
	FirstName         string  `bson:"firstName"`
	MiddleName        string  `bson:"middleName"`
	LastName          string  `bson:"lastName"`
	ProfilePictureUrl string  `bson:"profilePictureUrl"`
	Address           string  `bson:"address"`
	PhoneNumber       string  `bson:"phoneNumber"`
	EmailAddress      string  `bson:"emailAddress"`
	MaritalStatus     string  `bson:"maritalStatus"`
	Birthday          *string `bson:"birthday"`
	Gender            string  `bson:"gender"`
}

type personRepositoryImpl struct {
	ctx        context.Context
	db         *mongo.Database
	collection *mongo.Collection
}

// GetByEmail implements repositories.PersonRepository.
func (p *personRepositoryImpl) GetByEmail(email entities.EmailAddress) (*entities.Person, error) {
	var person Person
	err := p.collection.FindOne(p.ctx, bson.M{"emailAddress": string(email)}).Decode(&person)
	if(err != nil && err == mongo.ErrNoDocuments) {
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
	if(err != nil && err == mongo.ErrNoDocuments) {
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

func toPersonMongoModel(e entities.Person) Person {
	var birthday *string
	if e.Birthday != nil {
		str := fmt.Sprintf("%04d-%02d-%02d", e.Birthday.Year, e.Birthday.Month, e.Birthday.Day)
		birthday = &str
	} else {
		birthday = nil
	}

	return Person{
		ID:                e.ID,
		FirstName:         e.FirstName,
		MiddleName:        e.MiddleName,
		LastName:          e.LastName,
		ProfilePictureUrl: e.ProfilePictureUrl,
		Address:           e.Address,
		PhoneNumber:       string(e.PhoneNumber),
		EmailAddress:      string(e.EmailAddress),
		MaritalStatus:     e.MaritalStatus,
		Birthday:          birthday,
		Gender:            string(e.Gender),
	}
}

func toPersonEntities(p Person) entities.Person {
	var birthday *entities.YearMonthDay

	phones := make([]entities.PhoneNumber, len(p.PhoneNumber))

	for i, phone := range p.PhoneNumber {
		phones[i] = entities.PhoneNumber(phone)
	}

	if p.Birthday == nil {
		birthday = nil
	} else {
		birthday = &entities.YearMonthDay{}
		fmt.Sscanf(*p.Birthday, "%d-%d-%d", &birthday.Year, &birthday.Month, &birthday.Day)
	}

	return entities.Person{
		ID:                p.ID,
		FirstName:         p.FirstName,
		MiddleName:        p.MiddleName,
		LastName:          p.LastName,
		ProfilePictureUrl: p.ProfilePictureUrl,
		Address:           p.Address,
		PhoneNumber:       entities.PhoneNumber(p.PhoneNumber),
		EmailAddress:      entities.EmailAddress(p.EmailAddress),
		MaritalStatus:     p.MaritalStatus,
		Birthday:          birthday,
		Gender:            entities.Gender(p.Gender),
	}
}

func NewPersonRepository(ctx context.Context, db *mongo.Database) repositories.PersonRepository {
	return &personRepositoryImpl{
		ctx:        ctx,
		db:         db,
		collection: db.Collection("person"),
	}
}
