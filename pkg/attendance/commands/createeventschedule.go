package commands

import (
	"github.com/Goldwin/ies-pik-cms/pkg/attendance/entities"
	. "github.com/Goldwin/ies-pik-cms/pkg/common/commands"
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

	return CommandExecutionResult[entities.EventSchedule]{
		Status: ExecutionStatusFailed,
		Error: CommandErrorDetail{
			Code:    CreateEventScheduleCommandErrorAppError,
			Message: "Failed to Create Event Schedule",
		},
	}
}
