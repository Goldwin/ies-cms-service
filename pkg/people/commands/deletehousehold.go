package commands

import (
	"fmt"

	. "github.com/Goldwin/ies-pik-cms/pkg/common/commands"
	"github.com/Goldwin/ies-pik-cms/pkg/people/dto"
)

type DeleteHouseholdCommand struct {
	Input dto.HouseHoldInput
}

const (
	DeleteHouseholdErrorCodeHouseholdNotExistsError CommandErrorCode = 10401
	DeleteHouseholdErrorCodeDBError                 CommandErrorCode = 10402
)

func (cmd DeleteHouseholdCommand) Execute(ctx CommandContext) CommandExecutionResult[bool] {
	e, err := ctx.HouseholdRepository().Get(cmd.Input.ID)
	if err != nil {
		return CommandExecutionResult[bool]{
			Status: ExecutionStatusFailed,
			Error: CommandErrorDetail{
				Code:    DeleteHouseholdErrorCodeDBError,
				Message: fmt.Sprintf("Can't Delete New Household Info, Error: %s", err.Error()),
			},
		}
	}

	if e == nil {
		return CommandExecutionResult[bool]{
			Status: ExecutionStatusFailed,
			Error: CommandErrorDetail{
				Code:    DeleteHouseholdErrorCodeDBError,
				Message: fmt.Sprintf("Can't Delete Household Info, Error: Household with id %s does not exist", cmd.Input.ID),
			},
		}
	}

	err = ctx.HouseholdRepository().Delete(e)
	if(err != nil) {
		return CommandExecutionResult[bool]{
			Status: ExecutionStatusFailed,
			Error: CommandErrorDetail{
				Code:    DeleteHouseholdErrorCodeDBError,
				Message: fmt.Sprintf("Can't Delete Household Info, Error: %s", err.Error()),
			},
		}
	}

	return CommandExecutionResult[bool]{
		Status: ExecutionStatusSuccess,
		Result: true,
	}
}
