package commands

import (
	"github.com/Goldwin/ies-pik-cms/pkg/attendance/entities"
	. "github.com/Goldwin/ies-pik-cms/pkg/common/commands"
)

const (
	RemoveActivityIDDoesntExistsError       CommandErrorCode = 30211
	RemoveActivityScheduleDoesntExistsError CommandErrorCode = 30212
)

type RemoveScheduleActivityCommand struct {
	ScheduleID string
	ActivityID string
}

func (c RemoveScheduleActivityCommand) Execute(ctx CommandContext) CommandExecutionResult[entities.EventSchedule] {

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
				Code:    RemoveActivityScheduleDoesntExistsError,
				Message: "Schedule not found",
			},
		}
	}

	isActivityExists := false

	for i, activity := range schedule.Activities {
		if activity.ID == c.ActivityID {
			schedule.Activities = append(schedule.Activities[:i], schedule.Activities[i+1:]...)
			isActivityExists = true
			break
		}
	}

	if !isActivityExists {
		return CommandExecutionResult[entities.EventSchedule]{
			Status: ExecutionStatusFailed,
			Error: CommandErrorDetail{
				Code:    RemoveActivityIDDoesntExistsError,
				Message: "Activity not found",
			},
		}
	}

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
