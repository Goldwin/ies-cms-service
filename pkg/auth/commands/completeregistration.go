package commands

import (
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
)

type CompleteRegistrationCommand struct {
	FirstName  string
	MiddleName string
	LastName   string
	Email      string
}

func (cmd CompleteRegistrationCommand) Execute(ctx CommandContext) AppExecutionResult[dto.AuthData] {
	account, err := ctx.AccountRepository().GetAccount(entities.EmailAddress(cmd.Email))
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

	if cmd.FirstName == "" || cmd.LastName == "" {
		return AppExecutionResult[dto.AuthData]{
			Status: ExecutionStatusFailed,
			Error: AppErrorDetail{
				Code:    CompleteRegistrationErrorInvalidInput,
				Message: fmt.Sprintf("First Name and Last Name must be filled"),
			},
		}
	}

	account.Person = entities.Person{
		FirstName:  cmd.FirstName,
		MiddleName: cmd.MiddleName,
		LastName:   cmd.LastName,
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
