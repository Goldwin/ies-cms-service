package mongo

import (
	"context"

	"github.com/Goldwin/ies-pik-cms/pkg/attendance/commands"
	"github.com/Goldwin/ies-pik-cms/pkg/attendance/repositories"
	"go.mongodb.org/mongo-driver/mongo"
)

type commandContextImpl struct {
	ctx context.Context
	db  *mongo.Database

	eventRepository      repositories.EventRepository
	scheduleRepository   repositories.EventScheduleRepository
	attendanceRepository repositories.AttendanceRepository
}

// EventRepository implements commands.CommandContext.
func (c *commandContextImpl) EventRepository() repositories.EventRepository {
	if c.eventRepository == nil {
		c.eventRepository = NewEventRepository(c.ctx, c.db)
	}
	return c.eventRepository
}

// EventScheduleRepository implements commands.CommandContext.
func (c *commandContextImpl) EventScheduleRepository() repositories.EventScheduleRepository {
	if c.scheduleRepository == nil {
		c.scheduleRepository = NewEventScheduleRepository(c.ctx, c.db)
	}
	return c.scheduleRepository
}

func (c *commandContextImpl) AttendanceRepository() repositories.AttendanceRepository {
	if c.attendanceRepository == nil {
		c.attendanceRepository = NewAttendanceRepository(c.ctx, c.db)
	}
	return c.attendanceRepository
}

func NewCommandContext(ctx context.Context, db *mongo.Database) commands.CommandContext {
	return &commandContextImpl{
		ctx: ctx,
		db:  db,
	}
}
