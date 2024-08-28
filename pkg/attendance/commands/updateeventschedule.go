package commands

import (
	"time"

	"github.com/Goldwin/ies-pik-cms/pkg/attendance/entities"
	. "github.com/Goldwin/ies-pik-cms/pkg/common/commands"
)

const (
	UpdateEventScheduleErrorNotFound   CommandErrorCode = 30110
	UpdateEventScheduleValidationError CommandErrorCode = 30111
)

/*
Create Event Schedule Command.
This command will create new event schedule
*/
type UpdateEventScheduleCommand struct {
	ID             string
	Name           string
	ScheduleType   string
	TimezoneOffset int
	Days           []time.Weekday
	Date           time.Time
	StartDate      time.Time
	EndDate        time.Time
}

func (c UpdateEventScheduleCommand) Execute(ctx CommandContext) CommandExecutionResult[entities.EventSchedule] {
	eventSchedule, err := ctx.EventScheduleRepository().Get(c.ID)

	if err != nil {
		return CommandExecutionResult[entities.EventSchedule]{
			Status: ExecutionStatusFailed,
			Error:  CommandErrorDetailWorkerFailure(err),
		}
	}

	if eventSchedule == nil {
		return CommandExecutionResult[entities.EventSchedule]{
			Status: ExecutionStatusFailed,
			Error: CommandErrorDetail{
				Code:    UpdateEventScheduleValidationError,
				Message: "Event schedule not found",
			},
		}
	}

	eventSchedule.Name = c.Name
	eventSchedule.TimezoneOffset = c.TimezoneOffset
	eventSchedule.Type = entities.EventScheduleType(c.ScheduleType)

	if eventSchedule.Type == entities.EventScheduleTypeWeekly {
		eventSchedule.Days = c.Days
	} else if eventSchedule.Type == entities.EventScheduleTypeDaily {
		eventSchedule.StartDate = time.Date(c.StartDate.Year(), c.StartDate.Month(), c.StartDate.Day(), 0, 0, 0, 0, time.UTC)
		eventSchedule.EndDate = time.Date(c.EndDate.Year(), c.EndDate.Month(), c.EndDate.Day(), 0, 0, 0, 0, time.UTC)
	} else if eventSchedule.Type == entities.EventScheduleTypeOneTime {
		eventSchedule.Date = time.Date(c.Date.Year(), c.Date.Month(), c.Date.Day(), 0, 0, 0, 0, time.UTC)
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
