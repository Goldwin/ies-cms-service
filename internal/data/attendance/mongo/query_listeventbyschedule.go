package mongo

import (
	"context"

	"github.com/Goldwin/ies-pik-cms/pkg/attendance/dto"
	. "github.com/Goldwin/ies-pik-cms/pkg/attendance/queries"
	"github.com/Goldwin/ies-pik-cms/pkg/common/queries"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type listEventByScheduleImpl struct {
	ctx context.Context
	db  *mongo.Database
}

// Execute implements queries.ListEventBySchedule.
func (l *listEventByScheduleImpl) Execute(query ListEventByScheduleQuery) (ListEventByScheduleResult, queries.QueryErrorDetail) {
	cursor, err := l.db.Collection(EventCollection).Find(
		l.ctx,
		bson.M{"schedule_id": bson.M{"$eq": query.ScheduleID}, "date": bson.M{"$lt": query.LastDate}},
		options.Find().SetLimit(int64(query.Limit)),
		options.Find().SetSort(bson.D{{Key: "date", Value: -1}}),
	)

	if err != nil {
		return ListEventByScheduleResult{}, queries.QueryErrorDetail{
			Code:    500,
			Message: "Failed to connect to database",
		}
	}
	defer cursor.Close(l.ctx)
	var result = make([]dto.EventDTO, 0)
	for cursor.Next(l.ctx) {
		var eventModel EventModel
		if err := cursor.Decode(&eventModel); err != nil {
			return ListEventByScheduleResult{}, queries.QueryErrorDetail{
				Code:    500,
				Message: "Failed to Decode Event Information",
			}
		}
		result = append(result, dto.EventDTO{
			ID:         eventModel.ID,
			ScheduleID: eventModel.ScheduleID,
			Name:       eventModel.Name,
			Activities: lo.Map(eventModel.EventActivities, func(e EventActivityModel, _ int) dto.EventActivityDTO {
				return dto.EventActivityDTO{
					ID:   e.ID,
					Name: e.Name,
					Time: e.Time,
				}
			}),
			Date: eventModel.Date,
		})
	}
	return ListEventByScheduleResult{
		Data: result,
	}, queries.NoQueryError
}

func NewListEventBySchedule(ctx context.Context, db *mongo.Database) ListEventBySchedule {
	return &listEventByScheduleImpl{
		ctx: ctx,
		db:  db,
	}
}
