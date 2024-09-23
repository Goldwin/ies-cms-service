package commands

import (
	"fmt"
	"time"

	"github.com/Goldwin/ies-pik-cms/pkg/auth/dto"
	"github.com/Goldwin/ies-pik-cms/pkg/auth/entities"
	. "github.com/Goldwin/ies-pik-cms/pkg/common/commands"
)

const (
	ResetPasswordErrorInvalidInput          CommandErrorCode = 20403
	ResetPasswordErrorFailedToVerifyAccount CommandErrorCode = 20404
	ResetPasswordErrorTokenExpired          CommandErrorCode = 20405
)

type ResetPasswordCommand struct {
	Input dto.PasswordResetInput
}

func (cmd ResetPasswordCommand) Execute(ctx CommandContext) CommandExecutionResult[dto.PasswordResult] {

	code, err := ctx.PasswordResetCodeRepository().Get(cmd.Input.Email)

	if err != nil {
		return CommandExecutionResult[dto.PasswordResult]{
			Status: ExecutionStatusFailed,
			Error: CommandErrorDetail{
				Code:    SavePasswordErrorInvalidInput,
				Message: fmt.Sprintf("Failed to fetch Reset Token: %s", err.Error()),
			},
		}
	}

	token := code.Code

	if token != cmd.Input.Code {
		return CommandExecutionResult[dto.PasswordResult]{
			Status: ExecutionStatusFailed,
			Error: CommandErrorDetail{
				Code:    SavePasswordErrorInvalidInput,
				Message: fmt.Sprintf("Reset Token Mismatched"),
			},
		}
	}

	if code.ExpiryAt.Before(time.Now()) {
		return CommandExecutionResult[dto.PasswordResult]{
			Status: ExecutionStatusFailed,
			Error: CommandErrorDetail{
				Code:    SavePasswordErrorInvalidInput,
				Message: fmt.Sprintf("Reset Token Expired"),
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

	err = ctx.PasswordResetCodeRepository().Delete(code)

	return result
}
