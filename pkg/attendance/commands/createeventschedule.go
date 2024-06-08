package commands

import (
	"fmt"

	"github.com/Goldwin/ies-pik-cms/pkg/attendance/entities"
	. "github.com/Goldwin/ies-pik-cms/pkg/common/commands"
	"github.com/google/uuid"
)

const (
	CreateEventScheduleCommandErrorAppError CommandErrorCode = 30101
)

/*
Create Event Schedule Command.
This command will create new event schedule
*/
type CreateEventScheduleCommand struct {
	Name           string
	ScheduleType   string
	TimezoneOffset int
}

func (c CreateEventScheduleCommand) Execute(ctx CommandContext) CommandExecutionResult[entities.EventSchedule] {

	eventSchedule := &entities.EventSchedule{
		ID:             uuid.NewString(),
		Name:           c.Name,
		TimezoneOffset: c.TimezoneOffset,
		Type:           entities.EventScheduleType(c.ScheduleType),
	}

	result, err := ctx.EventScheduleRepository().Save(eventSchedule)

	if err != nil {

		return CommandExecutionResult[entities.EventSchedule]{
			Status: ExecutionStatusFailed,
			Error: CommandErrorDetail{
				Code:    CreateEventScheduleCommandErrorAppError,
				Message: fmt.Sprintf("Failed to Create Event Schedule: %s", err.Error()),
			},
		}
	}

	return CommandExecutionResult[entities.EventSchedule]{
		Status: ExecutionStatusSuccess,
		Result: *result,
	}
}
