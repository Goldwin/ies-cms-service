package events

import (
	"context"

	. "github.com/Goldwin/ies-pik-cms/pkg/common/commands"
	"github.com/Goldwin/ies-pik-cms/pkg/common/out"
	q "github.com/Goldwin/ies-pik-cms/pkg/common/queries"
	"github.com/Goldwin/ies-pik-cms/pkg/common/worker"
	"github.com/Goldwin/ies-pik-cms/pkg/events/commands"
	"github.com/Goldwin/ies-pik-cms/pkg/events/dto"
	"github.com/Goldwin/ies-pik-cms/pkg/events/queries"
)

type ChurchDataLayerComponent interface {
	CommandWorker() worker.UnitOfWork[commands.CommandContext]
	QueryWorker() worker.QueryWorker[queries.QueryContext]
}

type ChurchEventComponent interface {
	CheckIn(ctx context.Context, input dto.CheckInInput, output out.Output[[]dto.CheckInEvent])
	CreateEvent(ctx context.Context, input dto.ChurchEvent, output out.Output[dto.ChurchEvent])
	CreateSession(ctx context.Context, input dto.CreateSessionInput, output out.Output[dto.ChurchEventSession])
	SearchEvent(ctx context.Context, input queries.SearchEventQuery, output out.Output[queries.SearchEventResult])
}

type churchEventComponentImpl struct {
	commandWorker worker.UnitOfWork[commands.CommandContext]
	queryWorker   worker.QueryWorker[queries.QueryContext]
}

// SearchEvent implements ChurchEventComponent.
func (c *churchEventComponentImpl) SearchEvent(ctx context.Context, input queries.SearchEventQuery, output out.Output[queries.SearchEventResult]) {
	result, err := c.queryWorker.Query(ctx).SearchEvent().Execute(input)
	if err != q.NoQueryError {
		output.OnError(out.ConvertQueryErrorDetail(err))
	} else {
		output.OnSuccess(result)
	}
}

// CreateSession implements ChurchEventComponent.
func (c *churchEventComponentImpl) CreateSession(ctx context.Context, input dto.CreateSessionInput, output out.Output[dto.ChurchEventSession]) {
	var result CommandExecutionResult[dto.ChurchEventSession]
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
		output.OnError(out.ConvertCommandErrorDetail(result.Error))
	}
}

// CreateEvent implements ChurchEventComponent.
func (c *churchEventComponentImpl) CreateEvent(ctx context.Context, input dto.ChurchEvent, output out.Output[dto.ChurchEvent]) {
	var result CommandExecutionResult[dto.ChurchEvent]
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
		output.OnError(out.ConvertCommandErrorDetail(result.Error))
	}
}

// CheckIn implements ChurchEventComponent.
func (c *churchEventComponentImpl) CheckIn(ctx context.Context, input dto.CheckInInput, output out.Output[[]dto.CheckInEvent]) {
	var result CommandExecutionResult[[]dto.CheckInEvent]
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
		output.OnError(out.ConvertCommandErrorDetail(result.Error))
	}
}

func NewChurchEventComponent(datalayer ChurchDataLayerComponent) ChurchEventComponent {
	return &churchEventComponentImpl{
		commandWorker: datalayer.CommandWorker(),
		queryWorker:   datalayer.QueryWorker(),
	}
}
