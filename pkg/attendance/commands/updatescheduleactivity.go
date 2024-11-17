package commands

import (
	"github.com/Goldwin/ies-pik-cms/pkg/attendance/dto"
	"github.com/Goldwin/ies-pik-cms/pkg/attendance/entities"
	. "github.com/Goldwin/ies-pik-cms/pkg/common/commands"
	"github.com/samber/lo"
)

const (
	UpdateActivityIDDoesntExistsError       CommandErrorCode = 30211
	UpdateActivityScheduleDoesntExistsError CommandErrorCode = 30212
)

type UpdateEventScheduleActivityCommand struct {
	ScheduleID string
	ActivityID string
	Name       string
	Hour       int
	Minute     int
	Labels     []dto.ActivityLabelDTO
}

func (c UpdateEventScheduleActivityCommand) Execute(ctx CommandContext) CommandExecutionResult[entities.EventSchedule] {
	if c.Labels == nil {
		c.Labels = []dto.ActivityLabelDTO{}
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
				Code:    RemoveActivityScheduleDoesntExistsError,
				Message: "Schedule not found",
			},
		}
	}

	isActivityExists := false

	for i, activity := range schedule.Activities {
		if activity.ID == c.ActivityID {
			activity.Minute = c.Minute
			activity.Hour = c.Hour
			activity.Name = c.Name
			activity.Labels = lo.Map(c.Labels, func(label dto.ActivityLabelDTO, _ int) *entities.ActivityLabel {
				return label.ToEntity()
			})
			schedule.Activities[i] = activity
			isActivityExists = true
			break
		}
	}

	if !isActivityExists {
		return CommandExecutionResult[entities.EventSchedule]{
			Status: ExecutionStatusFailed,
			Error: CommandErrorDetail{
				Code:    UpdateActivityIDDoesntExistsError,
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
