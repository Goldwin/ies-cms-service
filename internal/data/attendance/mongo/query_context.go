package mongo

import (
	"context"

	"github.com/Goldwin/ies-pik-cms/pkg/attendance/queries"
	"go.mongodb.org/mongo-driver/mongo"
)

type queryContextImpl struct {
	ctx context.Context
	db  *mongo.Database

	getEvent       queries.GetEvent
	getSchedule    queries.GetEventSchedule
	listAttendance queries.ListEventAttendance
	listSchedules  queries.ListEventSchedule
	listBySchedule queries.ListEventBySchedule

	searchHousehold queries.SearchHousehold
}

// GetEventAttendanceSummary implements queries.QueryContext.
func (q *queryContextImpl) GetEventAttendanceSummary() queries.GetEventAttendanceSummary {
	return NewGetEventAttendanceSummary(q.ctx, q.db)
}

// SearchHousehold implements queries.QueryContext.
func (q *queryContextImpl) SearchHousehold() queries.SearchHousehold {
	return NewSearchHousehold(q.ctx, q.db)
}

// GetEvent implements queries.QueryContext.
func (q *queryContextImpl) GetEvent() queries.GetEvent {
	return NewGetEvent(q.ctx, q.db)
}

// GetEventSchedule implements queries.QueryContext.
func (q *queryContextImpl) GetEventSchedule() queries.GetEventSchedule {
	return NewGetEventSchedule(q.ctx, q.db)
}

// ListEventAttendance implements queries.QueryContext.
func (q *queryContextImpl) ListEventAttendance() queries.ListEventAttendance {
	return NewListEventAttendance(q.ctx, q.db)
}

// ListEventSchedules implements queries.QueryContext.
func (q *queryContextImpl) ListEventSchedules() queries.ListEventSchedule {
	return NewListEventSchedule(q.ctx, q.db)
}

// ListEventsBySchedule implements queries.QueryContext.
func (q *queryContextImpl) ListEventsBySchedule() queries.ListEventBySchedule {
	return NewListEventBySchedule(q.ctx, q.db)
}

func NewQueryContext(ctx context.Context, db *mongo.Database) queries.QueryContext {
	return &queryContextImpl{
		ctx: ctx,
		db:  db,
	}
}
