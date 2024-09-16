package mongo

import (
	"context"

	"github.com/Goldwin/ies-pik-cms/pkg/attendance/dto"
	. "github.com/Goldwin/ies-pik-cms/pkg/attendance/queries"
	"github.com/Goldwin/ies-pik-cms/pkg/common/queries"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type listEventByScheduleImpl struct {
	ctx context.Context
	db  *mongo.Database
}

// Execute implements queries.ListEventBySchedule.
func (l *listEventByScheduleImpl) Execute(query ListEventByScheduleFilter) (ListEventByScheduleResult, queries.QueryErrorDetail) {
	cursor, err := l.db.Collection(EventCollection).Find(
		l.ctx,
		bson.M{"scheduleId": bson.M{"$eq": query.ScheduleID}, "startDate": bson.M{"$gte": query.StartDate, "$lte": query.EndDate}},
		options.Find().SetLimit(int64(query.Limit)),
		options.Find().SetSort(bson.D{{Key: "_id", Value: -1}}),
	)

	if err != nil {
		return ListEventByScheduleResult{}, queries.InternalServerError(err)
	}
	defer cursor.Close(l.ctx)
	var result = make([]dto.EventDTO, 0)
	for cursor.Next(l.ctx) {
		var eventModel EventModel
		if err := cursor.Decode(&eventModel); err != nil {
			return ListEventByScheduleResult{}, queries.InternalServerError(err)
		}
		result = append(result, eventModel.ToDTO())
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
