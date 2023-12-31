package commands

import (
	. "github.com/Goldwin/ies-pik-cms/pkg/common/commands"
	"github.com/Goldwin/ies-pik-cms/pkg/events/dto"
)

type CreateEventSchedulesCommand struct {
}

func (cmd CreateEventSchedulesCommand) Execute() AppExecutionResult[dto.EventSchedule] {
	return AppExecutionResult[dto.EventSchedule]{
		Status: ExecutionStatusSuccess,
		Error:  AppErrorDetail{},
		Result: dto.EventSchedule{},
	}
}
