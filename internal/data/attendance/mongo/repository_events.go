package mongo

import (
	"context"

	"github.com/Goldwin/ies-pik-cms/pkg/attendance/entities"
	"github.com/Goldwin/ies-pik-cms/pkg/attendance/repositories"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type EventRepositoryImpl struct {
	ctx        context.Context
	db         *mongo.Database
	collection *mongo.Collection
}

// Delete implements repositories.EventRepository.
func (e *EventRepositoryImpl) Delete(event *entities.Event) error {
	_, err := e.collection.DeleteOne(e.ctx, bson.M{"_id": event.ID})
	return err
}

// Get implements repositories.EventRepository.
func (e *EventRepositoryImpl) Get(id string) (*entities.Event, error) {
	var model EventModel
	err := e.collection.FindOne(e.ctx, bson.M{"_id": id}).Decode(&model)

	if err == mongo.ErrNoDocuments {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}
	return model.ToEvent(), nil
}

// List implements repositories.EventRepository.
func (e *EventRepositoryImpl) List([]string) ([]*entities.Event, error) {
	var models []EventModel
	cursor, err := e.collection.Find(e.ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(e.ctx)
	if err = cursor.All(e.ctx, &models); err != nil {
		return nil, err
	}
	return lo.Map(models, func(model EventModel, _ int) *entities.Event {
		return model.ToEvent()
	}), nil
}

// Save implements repositories.EventRepository.
func (e *EventRepositoryImpl) Save(event *entities.Event) (*entities.Event, error) {

	model := toEventModel(event)
	_, err := e.collection.UpdateByID(e.ctx, event.ID, bson.M{"$set": model}, options.Update().SetUpsert(true))

	if err != nil {
		return nil, err
	}
	return event, nil
}

func NewEventRepository(ctx context.Context, db *mongo.Database) repositories.EventRepository {
	return &EventRepositoryImpl{
		ctx:        ctx,
		db:         db,
		collection: db.Collection(EventCollection),
	}
}
