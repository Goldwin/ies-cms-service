package mongo

import (
	"context"
	"time"

	"github.com/Goldwin/ies-pik-cms/pkg/events/entities"
	"github.com/Goldwin/ies-pik-cms/pkg/events/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type churchEventScheduleRepositoryImpl struct {
	ctx context.Context
	db  *mongo.Database
}

type ChurchEventSchedule struct {
	ID             string `bson:"_id"`
	Name           string `bson:"name"`
	DayOfWeek      int    `bson:"day_of_week"`
	Hours          int    `bson:"hours"`
	Minute         int    `bson:"minute"`
	TimezoneOffset int    `bson:"timezone_offset"`
}

// GetByTimezoneAndWeekDay implements repositories.ChurchEventScheduleRepository.
func (c *churchEventScheduleRepositoryImpl) GetByTimezoneAndWeekDay(timezoneOffset int, weekday time.Weekday) (*entities.ChurchEventSchedule, error) {
	var model ChurchEventSchedule
	err := c.db.Collection("schedules").FindOne(c.ctx, bson.M{"day_of_week": int(weekday), "timezone_offset": timezoneOffset}).Decode(&model)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &entities.ChurchEventSchedule{
		ID:             model.ID,
		Name:           model.Name,
		DayOfWeek:      model.DayOfWeek,
		Hours:          model.Hours,
		Minute:         model.Minute,
		TimezoneOffset: model.TimezoneOffset,
	}, nil
}

// Save implements repositories.ChurchEventScheduleRepository.
func (c *churchEventScheduleRepositoryImpl) Save(e entities.ChurchEventSchedule) error {
	_, err := c.db.Collection("schedules").InsertOne(c.ctx, e)
	return err
}

func NewChurchEventScheduleRepository(ctx context.Context, db *mongo.Database) repositories.ChurchEventScheduleRepository {
	return &churchEventScheduleRepositoryImpl{
		ctx: ctx,
		db:  db,
	}
}
