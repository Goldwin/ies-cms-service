package mongo

import (
	"context"

	. "github.com/Goldwin/ies-pik-cms/pkg/attendance/queries"
	"github.com/Goldwin/ies-pik-cms/pkg/common/queries"
	"go.mongodb.org/mongo-driver/mongo"
)

type listEventScheduleImpl struct {
	db *mongo.Database
}

// Execute implements queries.ListEventSchedule.
func (l *listEventScheduleImpl) Execute(query ListEventScheduleQuery) (ListEventScheduleResult, queries.QueryErrorDetail) {
	panic("unimplemented")
}

func NewListEventSchedule(ctx context.Context, db *mongo.Database) ListEventSchedule {
	return &listEventScheduleImpl{
		db: db,
	}
}
