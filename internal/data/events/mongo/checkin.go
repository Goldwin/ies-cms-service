package mongo

import (
	"context"
	"time"

	"github.com/Goldwin/ies-pik-cms/pkg/events/entities"
	"github.com/Goldwin/ies-pik-cms/pkg/events/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type eventCheckInRepositoryImpl struct {
	ctx context.Context
	db  *mongo.Database
}

type EventCheckInModel struct {
	ID        string      `bson:"_id"`
	Person    Person      `bson:"person"`
	Event     ChurchEvent `bson:"event"`
	CheckInAt time.Time   `bson:"check_in_at"`
}

// Get implements EventCheckInRepository.
func (e *eventCheckInRepositoryImpl) Get(id string) (*entities.CheckInEvent, error) {
	var model EventCheckInModel
	err := e.db.Collection("checkin").FindOne(e.ctx, bson.M{"_id": id}).Decode(&model)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}
	return &entities.CheckInEvent{
		ID: id,
		Person: entities.Person{
			ID:        id,
			FirstName: model.Person.FirstName,
			LastName:  model.Person.LastName,
		},
		Event: entities.ChurchEvent{
			ID:        id,
			Name:      model.Event.Name,
			StartTime: model.Event.StartTime,
		},
		CheckInAt: model.CheckInAt,
	}, err
}

// Save implements EventCheckInRepository.
func (e *eventCheckInRepositoryImpl) Save(checkIn entities.CheckInEvent) error {
	_, err := e.db.Collection("checkin").InsertOne(e.ctx, EventCheckInModel{
		ID:        checkIn.ID,
		Person:    Person{ID: checkIn.Person.ID, FirstName: checkIn.Person.FirstName, LastName: checkIn.Person.LastName},
		Event:     ChurchEvent{ID: checkIn.Event.ID, Name: checkIn.Event.Name, StartTime: checkIn.Event.StartTime},
		CheckInAt: checkIn.CheckInAt,
	})
	return err
}

func NewEventCheckInRepository(ctx context.Context, db *mongo.Database) repositories.EventCheckInRepository {
	return &eventCheckInRepositoryImpl{
		ctx: ctx,
		db:  db,
	}
}
