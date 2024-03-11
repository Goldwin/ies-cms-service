package commands

import (
	"bytes"
	"fmt"

	"github.com/Goldwin/ies-pik-cms/pkg/auth/dto"
	"github.com/Goldwin/ies-pik-cms/pkg/auth/entities"
	. "github.com/Goldwin/ies-pik-cms/pkg/common/commands"
)

const (
	CompleteRegistrationErrorAlreadyCompleted       CommandErrorCode = 20201
	CompleteRegistrationErrorUpdateFailure          CommandErrorCode = 20202
	CompleteRegistrationErrorAccountIsNotRegistered CommandErrorCode = 20203
	CompleteRegistrationErrorFailedToGetAccount     CommandErrorCode = 20203
	CompleteRegistrationErrorInvalidInput           CommandErrorCode = 20204
	CompleteRegistrationErrorPasswordMismatch       CommandErrorCode = 20205
)

type CompleteRegistrationCommand struct {
	Input dto.CompleteRegistrationInput
}

func (cmd CompleteRegistrationCommand) Execute(ctx CommandContext) CommandExecutionResult[dto.AuthData] {
	account, err := ctx.AccountRepository().GetAccount(entities.EmailAddress(cmd.Input.Email))
	if err != nil {
		return CommandExecutionResult[dto.AuthData]{
			Status: ExecutionStatusFailed,
			Error: CommandErrorDetail{
				Code:    CompleteRegistrationErrorFailedToGetAccount,
				Message: fmt.Sprintf("Failed to Get Account: %s", err.Error()),
			},
		}
	}
	if account == nil {
		return CommandExecutionResult[dto.AuthData]{
			Status: ExecutionStatusFailed,
			Error: CommandErrorDetail{
				Code:    CompleteRegistrationErrorAccountIsNotRegistered,
				Message: fmt.Sprintf("Account Is Not Registered"),
			},
		}
	}

	if len(account.Roles) > 0 {
		return CommandExecutionResult[dto.AuthData]{
			Status: ExecutionStatusFailed,
			Error: CommandErrorDetail{
				Code:    CompleteRegistrationErrorAlreadyCompleted,
				Message: fmt.Sprintf("Account Already Completed Registration"),
			},
		}
	}

	if cmd.Input.FirstName == "" || cmd.Input.LastName == "" {
		return CommandExecutionResult[dto.AuthData]{
			Status: ExecutionStatusFailed,
			Error: CommandErrorDetail{
				Code:    CompleteRegistrationErrorInvalidInput,
				Message: fmt.Sprintf("First Name and Last Name must be filled"),
			},
		}
	}

	if !bytes.Equal(cmd.Input.Password, cmd.Input.ConfirmPassword) {
		return CommandExecutionResult[dto.AuthData]{
			Status: ExecutionStatusFailed,
			Error: CommandErrorDetail{
				Code:    CompleteRegistrationErrorInvalidInput,
				Message: fmt.Sprintf("Password and Confirm Password should be same"),
			},
		}
	}

	account.Roles = []entities.Role{
		entities.ChurchMember,
	}

	account, err = ctx.AccountRepository().UpdateAccount(*account)

	if err != nil {
		return CommandExecutionResult[dto.AuthData]{
			Status: ExecutionStatusFailed,
			Error: CommandErrorDetail{
				Code:    CompleteRegistrationErrorUpdateFailure,
				Message: fmt.Sprintf("Failed to Update Account: %s", err.Error()),
			},
		}
	}

	// result := SavePasswordCommand{
	// 	Input: dto.PasswordInput{
	// 		Email:           cmd.Input.Email,
	// 		Password:        []byte(cmd.Input.Password),
	// 		ConfirmPassword: []byte(cmd.Input.ConfirmPassword),
	// 	},
	// }.Execute(ctx)

	// if result.Status != ExecutionStatusSuccess {
	// 	return CommandExecutionResult[dto.AuthData]{
	// 		Status: ExecutionStatusFailed,
	// 		Error:  result.Error,
	// 	}
	// }

	scopes := make([]string, 0)

	for _, role := range account.Roles {
		for _, scope := range role.Scopes {
			scopes = append(scopes, string(scope))
		}
	}

	return CommandExecutionResult[dto.AuthData]{
		Status: ExecutionStatusSuccess,
		Result: dto.AuthData{
			Email:  string(account.Email),
			Scopes: scopes,
		},
	}
}
