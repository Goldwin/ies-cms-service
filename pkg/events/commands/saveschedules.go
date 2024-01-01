package commands

import (
	"fmt"

	. "github.com/Goldwin/ies-pik-cms/pkg/common/commands"
	"github.com/Goldwin/ies-pik-cms/pkg/events/dto"
	"github.com/Goldwin/ies-pik-cms/pkg/events/entities"
	"github.com/Goldwin/ies-pik-cms/pkg/events/repositories"
)

const (
	CreateEventScheduleCommandsFailedToCreateEventSchedule AppErrorCode = 30201
)

type CreateEventScheduleCommand struct {
	Input dto.ChurchEventSchedule
}

func (cmd CreateEventScheduleCommand) Execute(ctx repositories.CommandContext) AppExecutionResult[dto.ChurchEventSchedule] {
	schedule := entities.ChurchEventSchedule{
		ID:             fmt.Sprintf("%d%2d%2d%+2d", cmd.Input.DayOfWeek, cmd.Input.Hours, cmd.Input.Minute, cmd.Input.TimezoneOffset),
		DayOfWeek:      cmd.Input.DayOfWeek,
		Hours:          cmd.Input.Hours,
		Minute:         cmd.Input.Minute,
		TimezoneOffset: cmd.Input.TimezoneOffset,
	}
	err := ctx.ChurchEventScheduleRepository().Save(schedule)

	if err != nil {
		return AppExecutionResult[dto.ChurchEventSchedule]{
			Status: ExecutionStatusFailed,
			Error: AppErrorDetail{
				Code:    CreateEventScheduleCommandsFailedToCreateEventSchedule,
				Message: fmt.Sprintf("Failed to Create Event Schedule: %s", err.Error()),
			},
		}
	}

	return AppExecutionResult[dto.ChurchEventSchedule]{
		Status: ExecutionStatusSuccess,
		Result: dto.ChurchEventSchedule{
			ID:             schedule.ID,
			DayOfWeek:      schedule.DayOfWeek,
			Hours:          schedule.Hours,
			Minute:         schedule.Minute,
			TimezoneOffset: schedule.TimezoneOffset,
		},
	}
}
