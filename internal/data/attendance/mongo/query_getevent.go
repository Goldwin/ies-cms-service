package mongo

import (
	"context"

	. "github.com/Goldwin/ies-pik-cms/pkg/attendance/queries"
	"github.com/Goldwin/ies-pik-cms/pkg/common/queries"
	"go.mongodb.org/mongo-driver/mongo"
)

type getEventImpl struct {
	db *mongo.Database
}

// Execute implements queries.GetEvent.
func (g *getEventImpl) Execute(query GetEventQuery) (GetEventResult, queries.QueryErrorDetail) {
	panic("unimplemented")
}

func NewGetEvent(ctx context.Context, db *mongo.Database) GetEvent {
	return &getEventImpl{
		db: db,
	}
}
