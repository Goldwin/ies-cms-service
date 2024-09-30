package commands

import (
	"fmt"

	"github.com/Goldwin/ies-pik-cms/pkg/auth/dto"
	"github.com/Goldwin/ies-pik-cms/pkg/auth/entities"
	. "github.com/Goldwin/ies-pik-cms/pkg/common/commands"
)

type InitializeRootAccountCommand struct {
	Email    string
	Password []byte
}

func (c InitializeRootAccountCommand) Execute(ctx CommandContext) CommandExecutionResult[dto.PasswordResult] {
	account, err := ctx.AccountRepository().Get(c.Email)

	if err != nil {
		return CommandExecutionResult[dto.PasswordResult]{
			Status: ExecutionStatusFailed,
			Error: CommandErrorDetail{
				Code:    SavePasswordErrorFailedToVerifyAccount,
				Message: fmt.Sprintf("Failed to Verify Account: %s", err.Error()),
			},
		}
	}

	if account != nil {
		return CommandExecutionResult[dto.PasswordResult]{
			Status: ExecutionStatusSuccess,
			Result: dto.PasswordResult{
				Email: c.Email,
			},
		}		
		
	}

	account = &entities.Account{
		Email: entities.EmailAddress(c.Email),
		Roles: []*entities.Role{&entities.Admin},
	}

	_, err = ctx.AccountRepository().Save(account)

	if err != nil {
		return CommandExecutionResult[dto.PasswordResult]{
			Status: ExecutionStatusFailed,
			Error: CommandErrorDetail{
				Code:    GrantAdminRoleErrorFailedToUpdateAccount,
				Message: fmt.Sprintf("Failed to Verify Account: %s", err.Error()),
			},
		}
	}

	return SavePasswordCommand{
		Email:    entities.EmailAddress(c.Email),
		Password: c.Password,
	}.Execute(ctx)	
}
