package mongo

import (
	"context"

	"github.com/Goldwin/ies-pik-cms/pkg/events/dto"
	"github.com/Goldwin/ies-pik-cms/pkg/events/queries"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type searchEventImpl struct {
	db  *mongo.Database
	ctx context.Context
}

// Execute implements queries.SearchEvent.
func (s *searchEventImpl) Execute(query queries.SearchEventQuery) (queries.SearchEventResult, error) {
	opts := options.Find().SetSort(bson.D{{Key: "_id", Value: 1}}).SetLimit(int64(query.Limit))
	cursor, err := s.db.Collection("event").Find(s.ctx, bson.M{"_id": bson.M{"$gt": query.LastID}}, opts)
	if err != nil {
		return queries.SearchEventResult{}, err
	}
	events := make([]dto.ChurchEvent, 0)
	defer cursor.Close(s.ctx)
	for cursor.Next(s.ctx) {
		var event ChurchEvent
		if err := cursor.Decode(&event); err != nil {
			return queries.SearchEventResult{}, err
		}
		events = append(events, toEventDTO(event))
	}

	return queries.SearchEventResult{}, nil
}

func toEventDTO(event ChurchEvent) dto.ChurchEvent {
	locations := make([]dto.Location, len(event.Locations))

	for i, l := range event.Locations {
		locations[i] = dto.Location{
			Name: l.Name,
			AgeFilter: dto.AgeFilter{
				From: l.AgeFilter.From,
				To:   l.AgeFilter.To,
			},
		}
	}

	return dto.ChurchEvent{
		ID:                     event.ID,
		Name:                   event.Name,
		Locations:              []dto.Location{},
		EventFrequency:         dto.Frequency(event.EventFrequency),
		LatestSessionStartTime: event.LatestSessionStartTime,
		ShowAt:                 event.LatestShowAt,
		HideAt:                 event.LatestHideAt,
	}
}

func SearchEvent(ctx context.Context, db *mongo.Database) queries.SearchEvent {
	return &searchEventImpl{
		db:  db,
		ctx: ctx,
	}
}
