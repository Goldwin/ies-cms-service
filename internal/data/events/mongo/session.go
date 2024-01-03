package mongo

import (
	"context"
	"time"

	"github.com/Goldwin/ies-pik-cms/pkg/events/entities"
	"github.com/Goldwin/ies-pik-cms/pkg/events/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ChurchEventSession struct {
	ID        string    `bson:"_id"`
	Name      string    `bson:"name"`
	SessionNo int       `bson:"session_no"`
	StartTime time.Time `bson:"start_time"`
	ShowAt    time.Time `bson:"show_at"`
	HideAt    time.Time `bson:"hide_at"`
}

type churchEventSessionRepositoryImpl struct {
	ctx context.Context
	db  mongo.Database
}

// Get implements repositories.ChurchEventSessionRepository.
func (c *churchEventSessionRepositoryImpl) Get(ID string) (*entities.ChurchEventSession, error) {
	var model ChurchEventSession
	err := c.db.Collection("sessions").FindOne(c.ctx, bson.M{"_id": ID}).Decode(&model)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &entities.ChurchEventSession{ID: model.ID, Name: model.Name, SessionNo: model.SessionNo, StartTime: model.StartTime, ShowAt: model.ShowAt, HideAt: model.HideAt}, nil
}

// Save implements repositories.ChurchEventSessionRepository.
func (c *churchEventSessionRepositoryImpl) Save(e entities.ChurchEventSession) error {
	_, err := c.db.Collection("sessions").UpdateByID(c.ctx, e.ID, bson.M{"$set": e}, options.Update().SetUpsert(true))
	return err
}

func NewChurchEventSessionRepository(ctx context.Context, db mongo.Database) repositories.ChurchEventSessionRepository {
	return &churchEventSessionRepositoryImpl{
		ctx: ctx,
		db:  db,
	}
}
