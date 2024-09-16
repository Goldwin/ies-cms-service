package mongo

import (
	"context"

	. "github.com/Goldwin/ies-pik-cms/pkg/attendance/queries"
	"github.com/Goldwin/ies-pik-cms/pkg/common/queries"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type getEventImpl struct {
	db  *mongo.Database
	ctx context.Context
}

// Execute implements queries.GetEvent.
func (g *getEventImpl) Execute(query GetEventFilter) (GetEventResult, queries.QueryErrorDetail) {
	var model EventModel
	err := g.db.Collection(EventCollection).FindOne(g.ctx, bson.M{"_id": query.EventID}).Decode(&model)

	if err == mongo.ErrNoDocuments {
		return GetEventResult{}, queries.ResourceNotFoundError("Event not found")
	}
	if err != nil {
		return GetEventResult{}, queries.InternalServerError(err)
	}

	return GetEventResult{
		Data: model.ToDTO(),
	}, queries.NoQueryError
}

func NewGetEvent(ctx context.Context, db *mongo.Database) GetEvent {
	return &getEventImpl{
		db:  db,
		ctx: ctx,
	}
}
