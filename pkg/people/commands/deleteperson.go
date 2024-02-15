package commands

import (
	"fmt"

	. "github.com/Goldwin/ies-pik-cms/pkg/common/commands"
	"github.com/Goldwin/ies-pik-cms/pkg/people/dto"
)

type DeletePersonCommand struct {
	Input dto.Person
}

const (
	DeletePersonErrorCodeUserNotExist CommandErrorCode = 10021
	DeletePersonErrorCodeDBError      CommandErrorCode = 10022
)

func (cmd DeletePersonCommand) Execute(ctx CommandContext) CommandExecutionResult[bool] {
	var err error
	c := cmd.Input

	personResult, err := ctx.PersonRepository().Get(c.ID)

	if personResult == nil {
		return CommandExecutionResult[bool]{
			Status: ExecutionStatusFailed,
			Error: CommandErrorDetail{
				Code:    DeletePersonErrorCodeUserNotExist,
				Message: fmt.Sprintf("Can't Delete Person Info, Error: Person with id %s does not exist", c.ID),
			},
		}
	}

	if err != nil {
		return CommandExecutionResult[bool]{
			Status: ExecutionStatusFailed,
			Error: CommandErrorDetail{
				Code:    DeletePersonErrorCodeDBError,
				Message: fmt.Sprintf("Can't Delete Person Info, Error: %s", err.Error()),
			},
		}
	}

	err = ctx.PersonRepository().DeletePerson(*personResult)

	if err != nil {
		return CommandExecutionResult[bool]{
			Status: ExecutionStatusFailed,
			Error: CommandErrorDetail{
				Code:    DeletePersonErrorCodeDBError,
				Message: fmt.Sprintf("Can't Delete Person Info, Error: %s", err.Error()),
			},
		}
	}

	return CommandExecutionResult[bool]{
		Status: ExecutionStatusSuccess,
		Error:  CommandErrorDetailNone,
		Result: true,
	}
}
