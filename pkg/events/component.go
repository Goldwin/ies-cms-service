package events

import (
	"context"
	"time"

	. "github.com/Goldwin/ies-pik-cms/pkg/common/commands"
	"github.com/Goldwin/ies-pik-cms/pkg/common/out"
	"github.com/Goldwin/ies-pik-cms/pkg/common/worker"
	"github.com/Goldwin/ies-pik-cms/pkg/events/commands"
	"github.com/Goldwin/ies-pik-cms/pkg/events/dto"
	"github.com/Goldwin/ies-pik-cms/pkg/events/repositories"
)

type ChurchDataLayerComponent interface {
	CommandWorker() worker.UnitOfWork[repositories.CommandContext]
}

type ChurchEventComponent interface {
	CheckIn(ctx context.Context, input dto.CheckInInput, output out.Output[dto.CheckInEvent])
	SaveEvent(ctx context.Context, input dto.ChurchEvent, output out.Output[dto.ChurchEvent])
	CreateEventSchedule(ctx context.Context, input dto.ChurchEventSchedule, output out.Output[dto.ChurchEventSchedule])
	CreateDailyEvent(ctx context.Context, input time.Weekday, output out.Output[[]dto.ChurchEvent])
}

type churchEventComponentImpl struct {
	commandWorker worker.UnitOfWork[repositories.CommandContext]
}

// CheckIn implements ChurchEventComponent.
func (c *churchEventComponentImpl) CheckIn(ctx context.Context, input dto.CheckInInput, output out.Output[dto.CheckInEvent]) {
	var result AppExecutionResult[dto.CheckInEvent]
	_ = c.commandWorker.Execute(ctx, func(ctx repositories.CommandContext) error {
		result = commands.CheckInCommands{
			Input: input,
		}.Execute(ctx)
		if result.Status == ExecutionStatusSuccess {
			go output.OnSuccess(result.Result)
		} else {
			go output.OnError(result.Error)
			return result.Error
		}
		return nil
	})
}

// CreateDailyEvent implements ChurchEventComponent.
func (c *churchEventComponentImpl) CreateDailyEvent(ctx context.Context, input time.Weekday, output out.Output[[]dto.ChurchEvent]) {
}

// SaveEvent implements ChurchEventComponent.
func (c *churchEventComponentImpl) SaveEvent(ctx context.Context, input dto.ChurchEvent, output out.Output[dto.ChurchEvent]) {

}

// CreateEventSchedule implements ChurchEventComponent.
func (c *churchEventComponentImpl) CreateEventSchedule(ctx context.Context, input dto.ChurchEventSchedule, output out.Output[dto.ChurchEventSchedule]) {
}

func NewChurchEventComponent(datalayer ChurchDataLayerComponent) ChurchEventComponent {
	return &churchEventComponentImpl{
		commandWorker: datalayer.CommandWorker(),
	}
}
