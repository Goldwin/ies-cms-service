package events

import (
	"context"

	. "github.com/Goldwin/ies-pik-cms/pkg/common/commands"
	"github.com/Goldwin/ies-pik-cms/pkg/common/out"
	"github.com/Goldwin/ies-pik-cms/pkg/common/worker"
	"github.com/Goldwin/ies-pik-cms/pkg/events/commands"
	"github.com/Goldwin/ies-pik-cms/pkg/events/dto"
)

type ChurchDataLayerComponent interface {
	CommandWorker() worker.UnitOfWork[commands.CommandContext]
}

type ChurchEventComponent interface {
	CheckIn(ctx context.Context, input dto.CheckInInput, output out.Output[[]dto.CheckInEvent])
	CreateEvent(ctx context.Context, input dto.ChurchEvent, output out.Output[dto.ChurchEvent])
	CreateSession(ctx context.Context, input dto.CreateSessionInput, output out.Output[dto.ChurchEventSession])
}

type churchEventComponentImpl struct {
	commandWorker worker.UnitOfWork[commands.CommandContext]
}

// CreateSession implements ChurchEventComponent.
func (c *churchEventComponentImpl) CreateSession(ctx context.Context, input dto.CreateSessionInput, output out.Output[dto.ChurchEventSession]) {
	var result AppExecutionResult[dto.ChurchEventSession]
	_ = c.commandWorker.Execute(ctx, func(ctx commands.CommandContext) error {
		result = commands.CreateChurchEventSessionCommand{
			EventID: input.EventID,
		}.Execute(ctx)
		if result.Status != ExecutionStatusSuccess {
			return result.Error
		}
		return nil
	})

	if result.Status == ExecutionStatusSuccess {
		output.OnSuccess(result.Result)
	} else {
		output.OnError(result.Error)
	}
}

// CreateEvent implements ChurchEventComponent.
func (c *churchEventComponentImpl) CreateEvent(ctx context.Context, input dto.ChurchEvent, output out.Output[dto.ChurchEvent]) {
	var result AppExecutionResult[dto.ChurchEvent]
	_ = c.commandWorker.Execute(ctx, func(ctx commands.CommandContext) error {
		result = commands.CreateEventCommands{
			Input: input,
		}.Execute(ctx)
		if result.Status != ExecutionStatusSuccess {
			return result.Error
		}
		return nil
	})
	if result.Status == ExecutionStatusSuccess {
		output.OnSuccess(result.Result)
	} else {
		output.OnError(result.Error)
	}
}

// CheckIn implements ChurchEventComponent.
func (c *churchEventComponentImpl) CheckIn(ctx context.Context, input dto.CheckInInput, output out.Output[[]dto.CheckInEvent]) {
	var result AppExecutionResult[[]dto.CheckInEvent]
	_ = c.commandWorker.Execute(ctx, func(ctx commands.CommandContext) error {
		result = commands.CheckInCommands{
			Input: input,
		}.Execute(ctx)
		if result.Status != ExecutionStatusSuccess {
			return result.Error
		}
		return nil
	})
	if result.Status == ExecutionStatusSuccess {
		output.OnSuccess(result.Result)
	} else {
		output.OnError(result.Error)
	}
}

func NewChurchEventComponent(datalayer ChurchDataLayerComponent) ChurchEventComponent {
	return &churchEventComponentImpl{
		commandWorker: datalayer.CommandWorker(),
	}
}
