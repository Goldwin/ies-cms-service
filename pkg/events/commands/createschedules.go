package commands

import (
	. "github.com/Goldwin/ies-pik-cms/pkg/common/commands"
	"github.com/Goldwin/ies-pik-cms/pkg/events/dto"
)

type CreateEventSchedulesCommand struct {
	input dto.ChurchEventSchedule
}

func (cmd CreateEventSchedulesCommand) Execute() AppExecutionResult[dto.ChurchEventSchedule] {
	return AppExecutionResult[dto.ChurchEventSchedule]{
		Status: ExecutionStatusSuccess,
		Error:  AppErrorDetail{},
		Result: dto.ChurchEventSchedule{},
	}
}
