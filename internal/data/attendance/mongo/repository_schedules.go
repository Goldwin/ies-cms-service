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

type eventScheduleRepositoryImpl struct {
	ctx        context.Context
	db         *mongo.Database
	collection *mongo.Collection
}

// Delete implements repositories.EventScheduleRepository.
func (e *eventScheduleRepositoryImpl) Delete(schedule *entities.EventSchedule) error {
	_, err := e.collection.DeleteOne(e.ctx, bson.M{"_id": schedule.ID})
	return err
}

// Get implements repositories.EventScheduleRepository.
func (e *eventScheduleRepositoryImpl) Get(id string) (*entities.EventSchedule, error) {
	var model EventScheduleModel
	err := e.collection.FindOne(e.ctx, bson.M{"_id": id}).Decode(&model)

	if err == mongo.ErrNoDocuments {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return model.ToEventSchedule(), nil
}

// List implements repositories.EventScheduleRepository.
func (e *eventScheduleRepositoryImpl) List(idList []string) ([]*entities.EventSchedule, error) {
	cursor, err := e.collection.Find(e.ctx, bson.M{"_id": bson.M{"$in": idList}})

	if err != nil {
		return nil, err
	}

	var models []EventScheduleModel
	if err = cursor.All(e.ctx, &models); err != nil {
		return nil, err
	}

	return lo.Map(models, func(model EventScheduleModel, _ int) *entities.EventSchedule {
		return model.ToEventSchedule()
	}), nil
}

// Save implements repositories.EventScheduleRepository.
func (e *eventScheduleRepositoryImpl) Save(schedule *entities.EventSchedule) (*entities.EventSchedule, error) {
	model := toEventScheduleModel(schedule)

	_, err := e.collection.UpdateByID(e.ctx, schedule.ID, bson.M{"$set": model}, options.Update().SetUpsert(true))
	if err != nil {
		return nil, err
	}
	return schedule, nil
}

func NewEventScheduleRepository(ctx context.Context, db *mongo.Database) repositories.EventScheduleRepository {
	return &eventScheduleRepositoryImpl{
		ctx:        ctx,
		db:         db,
		collection: db.Collection("eventSchedule"),
	}
}
