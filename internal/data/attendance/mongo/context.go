package mongo

import (
	"context"

	"github.com/Goldwin/ies-pik-cms/pkg/attendance/commands"
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
