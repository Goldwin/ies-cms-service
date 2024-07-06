package commands

import (
	"fmt"

	"github.com/Goldwin/ies-pik-cms/pkg/attendance/entities"
	. "github.com/Goldwin/ies-pik-cms/pkg/common/commands"
	"github.com/google/uuid"
)

const (
	AddActivityScheduleDoesntExistsError CommandErrorCode = 30201
	AddActivityInvalidTimeError          CommandErrorCode = 30202
	AddActivityLimitExceeded             CommandErrorCode = 30303
)

const (
	ActivityPerScheduleLimit int = 10
)

type AddEventScheduleActivityCommand struct {
	ScheduleID string
	Name       string
	Hour       int
	Minute     int
}

func (c AddEventScheduleActivityCommand) Execute(ctx CommandContext) CommandExecutionResult[entities.EventSchedule] {
	if c.Hour < 0 || c.Hour > 23 || c.Minute < 0 || c.Minute > 59 {
		return CommandExecutionResult[entities.EventSchedule]{
			Status: ExecutionStatusFailed,
			Error: CommandErrorDetail{
				Code:    AddActivityInvalidTimeError,
				Message: "Invalid hour or minute",
			},
		}
	}

	schedule, err := ctx.EventScheduleRepository().Get(c.ScheduleID)

	if err != nil {
		return CommandExecutionResult[entities.EventSchedule]{
			Status: ExecutionStatusFailed,
			Error:  CommandErrorDetailWorkerFailure(err),
		}
	}

	if schedule == nil {
		return CommandExecutionResult[entities.EventSchedule]{
			Status: ExecutionStatusFailed,
			Error: CommandErrorDetail{
				Code:    AddActivityScheduleDoesntExistsError,
				Message: "Schedule not found",
			},
		}
	}

	if (len(schedule.Activities) + 1) > ActivityPerScheduleLimit {
		return CommandExecutionResult[entities.EventSchedule]{
			Status: ExecutionStatusFailed,
			Error: CommandErrorDetail{
				Code:    AddActivityLimitExceeded,
				Message: fmt.Sprintf("Cannot add more than %d activities", ActivityPerScheduleLimit),
			},
		}
	}

	schedule.Activities = append(schedule.Activities, entities.EventScheduleActivity{
		ID:     uuid.NewString(),
		Name:   c.Name,
		Hour:   c.Hour,
		Minute: c.Minute,
	})

	result, err := ctx.EventScheduleRepository().Save(schedule)

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
