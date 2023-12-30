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
	ID                string    `bson:"_id"`
	FirstName         string    `bson:"firstName"`
	MiddleName        string    `bson:"middleName"`
	LastName          string    `bson:"lastName"`
	ProfilePictureUrl string    `bson:"profilePictureUrl"`
	Addresses         []Address `bson:"addresses"`
	PhoneNumbers      []string  `bson:"phoneNumbers"`
	EmailAddress      string    `bson:"emailAddress"`
	MaritalStatus     string    `bson:"maritalStatus"`
	Birthday          *string   `bson:"birthday"`
}

type Address struct {
	Line1      string `bson:"line1"`
	Line2      string `bson:"line2"`
	City       string `bson:"city"`
	Province   string `bson:"province"`
	PostalCode string `bson:"postalCode"`
}

type personRepositoryImpl struct {
	ctx        context.Context
	db         *mongo.Database
	collection *mongo.Collection
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
	addresses := make([]Address, len(e.Addresses))
	for i, address := range e.Addresses {
		addresses[i] = Address(address)
	}

	phones := make([]string, len(e.PhoneNumbers))
	for i, phone := range e.PhoneNumbers {
		phones[i] = string(phone)
	}

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
		Addresses:         addresses,
		PhoneNumbers:      phones,
		EmailAddress:      string(e.EmailAddress),
		MaritalStatus:     e.MaritalStatus,
		Birthday:          birthday,
	}
}

func toPersonEntities(p Person) entities.Person {
	var birthday *entities.YearMonthDay
	addresses := make([]entities.Address, len(p.Addresses))
	for i, address := range p.Addresses {
		addresses[i] = entities.Address(address)
	}

	phones := make([]entities.PhoneNumber, len(p.PhoneNumbers))

	for i, phone := range p.PhoneNumbers {
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
		Addresses:         addresses,
		PhoneNumbers:      phones,
		EmailAddress:      entities.EmailAddress(p.EmailAddress),
		MaritalStatus:     p.MaritalStatus,
		Birthday:          birthday,
	}
}

func NewPersonRepository(ctx context.Context, db *mongo.Database) repositories.PersonRepository {
	return &personRepositoryImpl{
		ctx:        ctx,
		db:         db,
		collection: db.Collection("person"),
	}
}
