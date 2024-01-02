package mongo

import (
	"context"
	"time"

	"github.com/Goldwin/ies-pik-cms/pkg/events/entities"
	"github.com/Goldwin/ies-pik-cms/pkg/events/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ChurchEvent struct {
	ID        string    `bson:"_id"`
	Name      string    `bson:"name"`
	StartTime time.Time `bson:"start_time"`
}

type churchEventRepositoryImpl struct {
	ctx context.Context
	db  *mongo.Database
}

// Get implements ChurchEventRepository.
func (c *churchEventRepositoryImpl) Get(id string) (*entities.ChurchEvent, error) {
	var model ChurchEvent
	err := c.db.Collection("events").FindOne(c.ctx, bson.M{"_id": id}).Decode(&model)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &entities.ChurchEvent{
		ID:        id,
		Name:      model.ID,
		StartTime: model.StartTime.UTC(),
	}, nil
}

// Save implements ChurchEventRepository.
func (c *churchEventRepositoryImpl) Save(e entities.ChurchEvent) error {
	_, err := c.db.Collection("events").InsertOne(c.ctx, ChurchEvent{
		ID:        e.ID,
		Name:      e.Name,
		StartTime: e.StartTime.UTC(),
	})
	return err
}

func NewChurchEventRepository(ctx context.Context, db *mongo.Database) repositories.ChurchEventRepository {
	return &churchEventRepositoryImpl{
		ctx: ctx,
		db:  db,
	}
}
