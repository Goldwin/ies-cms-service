package commands

import (
	"time"

	"github.com/Goldwin/ies-pik-cms/pkg/attendance/entities"
	. "github.com/Goldwin/ies-pik-cms/pkg/common/commands"
	"github.com/google/uuid"
)

const (
	CreateEventCommandsErrorScheduleDoesntExists CommandErrorCode = 30301
	CreateEventCommandsErrorInvalidTargetDate    CommandErrorCode = 30302
)

type CreateEventCommand struct {
	ScheduleID string
	TargetDate time.Time
}

func (c CreateEventCommand) Execute(ctx CommandContext) CommandExecutionResult[entities.Event] {
	schedule, err := ctx.EventScheduleRepository().Get(c.ScheduleID)

	if err != nil {
		return CommandExecutionResult[entities.Event]{
			Status: ExecutionStatusFailed,
			Error:  CommandErrorDetailWorkerFailure(err),
		}
	}

	if schedule == nil {
		return CommandExecutionResult[entities.Event]{
			Status: ExecutionStatusFailed,
			Error: CommandErrorDetail{
				Code:    CreateEventCommandsErrorScheduleDoesntExists,
				Message: "Schedule not found",
			},
		}
	}

	nextEvent, err := schedule.CreateNextEvent(c.TargetDate)

	if err != nil {
		return CommandExecutionResult[entities.Event]{
			Status: ExecutionStatusFailed,
			Error: CommandErrorDetail{
				Code:    CreateEventCommandsErrorInvalidTargetDate,
				Message: err.Error(),
			},
		}
	}

	nextEvent.ID = uuid.NewString()

	nextEvent, err = ctx.EventRepository().Save(nextEvent)

	return CommandExecutionResult[entities.Event]{
		Status: ExecutionStatusSuccess,
		Result: *nextEvent,
	}
}
