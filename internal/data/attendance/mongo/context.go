package mongo

import (
	"context"

	"github.com/Goldwin/ies-pik-cms/pkg/attendance/commands"
	"github.com/Goldwin/ies-pik-cms/pkg/attendance/queries"
	"github.com/Goldwin/ies-pik-cms/pkg/attendance/repositories"
	"go.mongodb.org/mongo-driver/mongo"
)

type commandContextImpl struct {
	eventRepository    repositories.EventRepository
	scheduleRepository repositories.EventScheduleRepository
}

// EventRepository implements commands.CommandContext.
func (c *commandContextImpl) EventRepository() repositories.EventRepository {
	return c.eventRepository
}

// EventScheduleRepository implements commands.CommandContext.
func (c *commandContextImpl) EventScheduleRepository() repositories.EventScheduleRepository {
	return c.scheduleRepository
}

func NewCommandContext(ctx context.Context, db *mongo.Database) commands.CommandContext {
	return &commandContextImpl{
		eventRepository:    NewEventRepository(ctx, db),
		scheduleRepository: NewEventScheduleRepository(ctx, db),
	}
}

type queryContextImpl struct {
	ctx context.Context
	db  *mongo.Database

	getEvent    queries.GetEvent
	getSchedule queries.GetEventSchedule
	listAttendance queries.ListEventAttendance
	listSchedules  queries.ListEventSchedule
	listBySchedule queries.ListEventBySchedule
}

// GetEvent implements queries.QueryContext.
func (q *queryContextImpl) GetEvent() queries.GetEvent {
	panic("unimplemented")
}

// GetEventSchedule implements queries.QueryContext.
func (q *queryContextImpl) GetEventSchedule() queries.GetEventSchedule {
	panic("unimplemented")
}

// ListEventAttendance implements queries.QueryContext.
func (q *queryContextImpl) ListEventAttendance() queries.ListEventAttendance {
	panic("unimplemented")
}

// ListEventSchedules implements queries.QueryContext.
func (q *queryContextImpl) ListEventSchedules() queries.ListEventSchedule {
	panic("unimplemented")
}

// ListEventsBySchedule implements queries.QueryContext.
func (q *queryContextImpl) ListEventsBySchedule() queries.ListEventBySchedule {
	panic("unimplemented")
}

func NewQueryContext(ctx context.Context, db *mongo.Database) queries.QueryContext {
	return &queryContextImpl{
	}
}
