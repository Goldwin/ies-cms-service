package commands

import (
	"bytes"
	"fmt"

	"github.com/Goldwin/ies-pik-cms/pkg/auth/dto"
	"github.com/Goldwin/ies-pik-cms/pkg/auth/entities"
	. "github.com/Goldwin/ies-pik-cms/pkg/common/commands"
)

const (
	CompleteRegistrationErrorAlreadyCompleted       AppErrorCode = 20201
	CompleteRegistrationErrorUpdateFailure          AppErrorCode = 20202
	CompleteRegistrationErrorAccountIsNotRegistered AppErrorCode = 20203
	CompleteRegistrationErrorFailedToGetAccount     AppErrorCode = 20203
	CompleteRegistrationErrorInvalidInput           AppErrorCode = 20204
	CompleteRegistrationErrorPasswordMismatch       AppErrorCode = 20205
)

type CompleteRegistrationCommand struct {
	Input dto.CompleteRegistrationInput
}

func (cmd CompleteRegistrationCommand) Execute(ctx CommandContext) AppExecutionResult[dto.AuthData] {
	account, err := ctx.AccountRepository().GetAccount(entities.EmailAddress(cmd.Input.Email))
	if err != nil {
		return AppExecutionResult[dto.AuthData]{
			Status: ExecutionStatusFailed,
			Error: AppErrorDetail{
				Code:    CompleteRegistrationErrorFailedToGetAccount,
				Message: fmt.Sprintf("Failed to Get Account: %s", err.Error()),
			},
		}
	}
	if account == nil {
		return AppExecutionResult[dto.AuthData]{
			Status: ExecutionStatusFailed,
			Error: AppErrorDetail{
				Code:    CompleteRegistrationErrorAccountIsNotRegistered,
				Message: fmt.Sprintf("Account Is Not Registered"),
			},
		}
	}

	if len(account.Roles) > 0 {
		return AppExecutionResult[dto.AuthData]{
			Status: ExecutionStatusFailed,
			Error: AppErrorDetail{
				Code:    CompleteRegistrationErrorAlreadyCompleted,
				Message: fmt.Sprintf("Account Already Completed Registration"),
			},
		}
	}

	if cmd.Input.FirstName == "" || cmd.Input.LastName == "" {
		return AppExecutionResult[dto.AuthData]{
			Status: ExecutionStatusFailed,
			Error: AppErrorDetail{
				Code:    CompleteRegistrationErrorInvalidInput,
				Message: fmt.Sprintf("First Name and Last Name must be filled"),
			},
		}
	}

	if !bytes.Equal(cmd.Input.Password, cmd.Input.ConfirmPassword) {
		return AppExecutionResult[dto.AuthData]{
			Status: ExecutionStatusFailed,
			Error: AppErrorDetail{
				Code:    CompleteRegistrationErrorInvalidInput,
				Message: fmt.Sprintf("Password and Confirm Password should be same"),
			},
		}
	}

	account.Person = entities.Person{
		FirstName:  cmd.Input.FirstName,
		MiddleName: cmd.Input.MiddleName,
		LastName:   cmd.Input.LastName,
	}
	account.Roles = []entities.Role{
		entities.ChurchMember,
	}

	account, err = ctx.AccountRepository().UpdateAccount(*account)

	if err != nil {
		return AppExecutionResult[dto.AuthData]{
			Status: ExecutionStatusFailed,
			Error: AppErrorDetail{
				Code:    CompleteRegistrationErrorUpdateFailure,
				Message: fmt.Sprintf("Failed to Update Account: %s", err.Error()),
			},
		}
	}

	result := SavePasswordCommand{
		Input: dto.PasswordInput{
			Email:           cmd.Input.Email,
			Password:        []byte(cmd.Input.Password),
			ConfirmPassword: []byte(cmd.Input.ConfirmPassword),
		},
	}.Execute(ctx)

	if result.Status != ExecutionStatusSuccess {
		return AppExecutionResult[dto.AuthData]{
			Status: ExecutionStatusFailed,
			Error:  result.Error,
		}
	}

	scopes := make([]string, 0)

	for _, role := range account.Roles {
		for _, scope := range role.Scopes {
			scopes = append(scopes, string(scope))
		}
	}

	return AppExecutionResult[dto.AuthData]{
		Status: ExecutionStatusSuccess,
		Result: dto.AuthData{
			ID:         account.Person.ID,
			FirstName:  account.Person.FirstName,
			MiddleName: account.Person.MiddleName,
			LastName:   account.Person.LastName,
			Email:      string(account.Email),
			Scopes:     scopes,
		},
	}
}
