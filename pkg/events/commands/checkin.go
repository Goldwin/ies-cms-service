package commands

import (
	. "github.com/Goldwin/ies-pik-cms/pkg/common/commands"
	"github.com/Goldwin/ies-pik-cms/pkg/events/dto"
)

type CheckInCommands struct {
}

func (cmd CheckInCommands) Execute() AppExecutionResult[dto.ChurchEvent] {

	return AppExecutionResult[dto.ChurchEvent]{
		Status: ExecutionStatusSuccess,
		Result: dto.ChurchEvent{
			ID:     "",
			Year:   0,
			Month:  0,
			Day:    0,
			Hours:  0,
			Minute: 0,
		},
	}
}
