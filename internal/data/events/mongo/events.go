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

type AgeFilter struct {
	From int `bson:"from"`
	To   int `bson:"to"`
}

type GenderFilter string

type Location struct {
	Name      string    `bson:"name"`
	AgeFilter AgeFilter `bson:"age_filter"`
}

type ChurchEvent struct {
	ID                     string     `bson:"_id"`
	Name                   string     `bson:"name"`
	Locations              []Location `bson:"locations"`
	EventFrequency         string     `bson:"event_frequency"`
	LatestSessionStartTime time.Time  `bson:"latest_session_start_time"`
	LatestShowAt           time.Time  `bson:"latest_show_at"`
	LatestHideAt           time.Time  `bson:"latest_hide_at"`
	LatestSessionNo        int        `bson:"latest_session_no"`
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

	locations := []entities.Location{}

	for _, l := range model.Locations {
		locations = append(locations, entities.Location{
			Name: l.Name,
			AgeFilter: entities.AgeFilter{
				From: l.AgeFilter.From,
				To:   l.AgeFilter.To,
			},
		})
	}

	return &entities.ChurchEvent{
		ID:                     id,
		Name:                   model.ID,
		Locations:              []entities.Location{},
		EventFrequency:         entities.Frequency(model.EventFrequency),
		LatestSessionStartTime: model.LatestSessionStartTime,
		LatestShowAt:           model.LatestShowAt,
		LatestHideAt:           model.LatestHideAt,
		LatestSessionNo:        model.LatestSessionNo,
	}, nil
}

// Save implements ChurchEventRepository.
func (c *churchEventRepositoryImpl) Save(e entities.ChurchEvent) error {
	locations := []Location{}

	for _, l := range e.Locations {
		locations = append(locations, Location{
			Name:      l.Name,
			AgeFilter: AgeFilter{From: l.AgeFilter.From, To: l.AgeFilter.To},
		})
	}
	_, err := c.db.Collection("events").UpdateByID(c.ctx, e.ID, bson.M{
		"$set": bson.M{
			"name":                    e.Name,
			"locations":               locations,
			"event_frequency":         string(e.EventFrequency),
			"latest_session_start_at": e.LatestSessionStartTime,
			"latest_show_at":          e.LatestShowAt,
			"latest_hide_at":          e.LatestHideAt,
			"latest_session_no":       e.LatestSessionNo,
		},
	}, options.Update().SetUpsert(true))
	return err
}

func NewChurchEventRepository(ctx context.Context, db *mongo.Database) repositories.ChurchEventRepository {
	return &churchEventRepositoryImpl{
		ctx: ctx,
		db:  db,
	}
}
