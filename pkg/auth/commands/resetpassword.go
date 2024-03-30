package commands

import (
	"fmt"

	"github.com/Goldwin/ies-pik-cms/pkg/auth/dto"
	"github.com/Goldwin/ies-pik-cms/pkg/auth/entities"
	. "github.com/Goldwin/ies-pik-cms/pkg/common/commands"
)

const (
	ResetPasswordErrorInvalidInput          CommandErrorCode = 20401
	ResetPasswordErrorFailedToVerifyAccount CommandErrorCode = 20402
)

type ResetPasswordCommand struct {
	Input dto.PasswordResetInput
}

func (cmd ResetPasswordCommand) Execute(ctx CommandContext) CommandExecutionResult[dto.PasswordResult] {

	token, err := ctx.PasswordRepository().GetResetCode(entities.EmailAddress(cmd.Input.Email))

	if err != nil {
		return CommandExecutionResult[dto.PasswordResult]{
			Status: ExecutionStatusFailed,
			Error: CommandErrorDetail{
				Code:    SavePasswordErrorInvalidInput,
				Message: fmt.Sprintf("Failed to fetch Reset Token: %s", err.Error()),
			},
		}
	}

	if token != cmd.Input.Code {
		return CommandExecutionResult[dto.PasswordResult]{
			Status: ExecutionStatusFailed,
			Error: CommandErrorDetail{
				Code:    SavePasswordErrorInvalidInput,
				Message: fmt.Sprintf("Reset Token Mismatched"),
			},
		}
	}

	savePasswordCmd := SavePasswordCommand{
		Email:    entities.EmailAddress(cmd.Input.Email),
		Password: []byte(cmd.Input.Password),
	}

	result := savePasswordCmd.Execute(ctx)

	if result.Status != ExecutionStatusSuccess {
		return result
	}

	err = ctx.PasswordRepository().DeleteResetToken(entities.EmailAddress(cmd.Input.Email))

	return result
}
