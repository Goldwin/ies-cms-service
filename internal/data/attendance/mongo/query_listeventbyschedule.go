package mongo

import (
	"context"

	. "github.com/Goldwin/ies-pik-cms/pkg/attendance/queries"
	"github.com/Goldwin/ies-pik-cms/pkg/common/queries"
	"go.mongodb.org/mongo-driver/mongo"
)

type listEventByScheduleImpl struct {
	db *mongo.Database
}

// Execute implements queries.ListEventBySchedule.
func (l *listEventByScheduleImpl) Execute(query ListEventByScheduleQuery) (ListEventByScheduleResult, queries.QueryErrorDetail) {
	panic("unimplemented")
}

func NewListEventBySchedule(ctx context.Context, db *mongo.Database) ListEventBySchedule {
	return &listEventByScheduleImpl{
		db: db,
	}
}
