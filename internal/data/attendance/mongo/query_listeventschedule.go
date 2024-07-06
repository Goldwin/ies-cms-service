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

type listEventScheduleImpl struct {
	ctx context.Context
	db  *mongo.Database
}

// Execute implements queries.ListEventSchedule.
func (l *listEventScheduleImpl) Execute(query ListEventScheduleQuery) (ListEventScheduleResult, queries.QueryErrorDetail) {
	cursor, err := l.db.Collection(EventScheduleCollection).Find(l.ctx,
		bson.M{
			"_id": bson.M{"$gt": query.LastID},
		},
		options.Find().SetLimit(int64(query.Limit)),
	)

	if err != nil {
		return ListEventScheduleResult{}, queries.InternalServerError(err)
	}
	defer cursor.Close(l.ctx)
	var result = make([]dto.EventScheduleDTO, 0)
	for cursor.Next(l.ctx) {
		var eventScheduleModel EventScheduleModel
		if err := cursor.Decode(&eventScheduleModel); err != nil {
			return ListEventScheduleResult{}, queries.InternalServerError(err)
		}
		result = append(result, dto.EventScheduleDTO{
			ID:             eventScheduleModel.ID,
			Name:           eventScheduleModel.Name,
			Type:           eventScheduleModel.Type,
			TimezoneOffset: eventScheduleModel.TimezoneOffset,
			Days:           eventScheduleModel.Days,
			Date:           eventScheduleModel.Date,
			StartDate:      eventScheduleModel.StartDate,
			EndDate:        eventScheduleModel.EndDate,
			Activities: lo.Map(
				eventScheduleModel.Activities, func(e EventScheduleActivityModel, _ int) dto.EventScheduleActivityDTO {
					return dto.EventScheduleActivityDTO{
						ScheduleID: eventScheduleModel.ID,
						ID:         e.ID,
						Name:       e.Name,
						Hour:       e.Hour,
						Minute:     e.Minute,
					}
				},
			),
		})
	}
	return ListEventScheduleResult{Data: result}, queries.NoQueryError

}

func NewListEventSchedule(ctx context.Context, db *mongo.Database) ListEventSchedule {
	return &listEventScheduleImpl{
		ctx: ctx,
		db:  db,
	}
}
