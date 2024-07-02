package mongo

import (
	"context"

	. "github.com/Goldwin/ies-pik-cms/pkg/attendance/queries"
	"github.com/Goldwin/ies-pik-cms/pkg/common/queries"
	"go.mongodb.org/mongo-driver/mongo"
)

type listEventAttendanceImpl struct {
	db *mongo.Database
}

// Execute implements queries.ListEventAttendance.
func (l *listEventAttendanceImpl) Execute(query ListEventAttendanceQuery) (ListEventAttendanceResult, queries.QueryErrorDetail) {
	panic("unimplemented")
}

func NewListEventAttendance(ctx context.Context, db *mongo.Database) ListEventAttendance {
	return &listEventAttendanceImpl{
		db: db,
	}
}
