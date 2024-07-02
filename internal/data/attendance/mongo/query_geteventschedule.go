package mongo

import (
	"context"

	. "github.com/Goldwin/ies-pik-cms/pkg/attendance/queries"
	"github.com/Goldwin/ies-pik-cms/pkg/common/queries"
	"go.mongodb.org/mongo-driver/mongo"
)

type getEventScheduleImpl struct {
	db *mongo.Database
}

// Execute implements queries.GetEventSchedule.
func (g *getEventScheduleImpl) Execute(query GetEventScheduleQuery) (GetEventScheduleResult, queries.QueryErrorDetail) {
	panic("unimplemented")
}

func NewGetEventSchedule(ctx context.Context, db *mongo.Database) GetEventSchedule {
	return &getEventScheduleImpl{
		db: db,
	}
}
