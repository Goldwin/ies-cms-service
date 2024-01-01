package commands

import (
	. "github.com/Goldwin/ies-pik-cms/pkg/common/commands"
	"github.com/Goldwin/ies-pik-cms/pkg/events/dto"
)

type CreateEventCommand struct {
	Input dto.CreateEventInput
}

func (cmd CreateEventCommand) Execute() AppExecutionResult[dto.ChurchEvent] {

	return AppExecutionResult[dto.ChurchEvent]{
		Status: ExecutionStatusSuccess,
		Result: dto.ChurchEvent{
			ID: "",
		},
	}
}
