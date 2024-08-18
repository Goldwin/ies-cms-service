package mongo

import (
	"context"

	"github.com/Goldwin/ies-pik-cms/pkg/attendance/dto"
	. "github.com/Goldwin/ies-pik-cms/pkg/attendance/queries"
	"github.com/Goldwin/ies-pik-cms/pkg/common/queries"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type getEventScheduleImpl struct {
	db  *mongo.Database
	ctx context.Context
}

// Execute implements queries.GetEventSchedule.
func (g *getEventScheduleImpl) Execute(query GetEventScheduleQuery) (GetEventScheduleResult, queries.QueryErrorDetail) {
	var model EventScheduleModel
	err := g.db.Collection(EventScheduleCollection).FindOne(g.ctx, bson.M{"_id": query.ScheduleID}).Decode(&model)

	if err == mongo.ErrNoDocuments {
		return GetEventScheduleResult{}, queries.ResourceNotFoundError("Event schedule not found")
	}

	if err != nil {
		return GetEventScheduleResult{}, queries.InternalServerError(err)
	}

	return GetEventScheduleResult{
		Data: dto.EventScheduleDTO{
			ID:             model.ID,
			Name:           model.Name,
			TimezoneOffset: model.TimezoneOffset,
			Type:           model.Type,
			Activities: lo.Map(model.Activities, func(e EventScheduleActivityModel, _ int) dto.EventScheduleActivityDTO {
				return dto.EventScheduleActivityDTO{
					ID:         e.ID,
					ScheduleID: model.ID,
					Name:       e.Name,
					Hour:       e.Hour,
					Minute:     e.Minute,
				}
			}),
			Date:      model.Date,
			Days:      model.Days,
			StartDate: model.StartDate,
			EndDate:   model.EndDate,
		},
	}, queries.NoQueryError
}

func NewGetEventSchedule(ctx context.Context, db *mongo.Database) GetEventSchedule {
	return &getEventScheduleImpl{
		db:  db,
		ctx: ctx,
	}
}
