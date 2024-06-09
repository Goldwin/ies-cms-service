package commands

import (
	"time"

	"github.com/Goldwin/ies-pik-cms/pkg/attendance/entities"
	. "github.com/Goldwin/ies-pik-cms/pkg/common/commands"
	"github.com/google/uuid"
)

const (
	CreateEventScheduleValidationError CommandErrorCode = 30101
)

/*
Create Event Schedule Command.
This command will create new event schedule
*/
type CreateEventScheduleCommand struct {
	Name           string
	ScheduleType   string
	TimezoneOffset int
	Days           []time.Weekday
	Date           time.Time
	StartDate      time.Time
	EndDate        time.Time
}

func (c CreateEventScheduleCommand) Execute(ctx CommandContext) CommandExecutionResult[entities.EventSchedule] {

	eventSchedule := &entities.EventSchedule{
		ID:             uuid.NewString(),
		Name:           c.Name,
		TimezoneOffset: c.TimezoneOffset,
		Type:           entities.EventScheduleType(c.ScheduleType),
	}

	validationMsg := eventSchedule.IsValid()

	if validationMsg != "" {
		return CommandExecutionResult[entities.EventSchedule]{
			Status: ExecutionStatusFailed,
			Error: CommandErrorDetail{
				Code:    CreateEventScheduleValidationError,
				Message: validationMsg,
			},
		}
	}

	result, err := ctx.EventScheduleRepository().Save(eventSchedule)

	if err != nil {

		return CommandExecutionResult[entities.EventSchedule]{
			Status: ExecutionStatusFailed,
			Error:  CommandErrorDetailWorkerFailure(err),
		}
	}

	return CommandExecutionResult[entities.EventSchedule]{
		Status: ExecutionStatusSuccess,
		Result: *result,
	}
}
