package commands

import (
	"fmt"

	"github.com/Goldwin/ies-pik-cms/pkg/auth/dto"
	"github.com/Goldwin/ies-pik-cms/pkg/auth/entities"
	. "github.com/Goldwin/ies-pik-cms/pkg/common/commands"
)

const (
	GrantAdminRoleErrorFailedToVerifyAccount CommandErrorCode = 20501
	GrantAdminRoleErrorFailedToUpdateAccount CommandErrorCode = 20502
)

// A command to grant admin role to user
type GrantAdminRoleCommand struct {
	Email string
}

func (cmd GrantAdminRoleCommand) Execute(ctx CommandContext) CommandExecutionResult[dto.AuthData] {
	account, err := ctx.AccountRepository().Get(cmd.Email)
	if err != nil {
		return CommandExecutionResult[dto.AuthData]{
			Status: ExecutionStatusFailed,
			Error: CommandErrorDetail{
				Code:    GrantAdminRoleErrorFailedToVerifyAccount,
				Message: fmt.Sprintf("Failed to Verify Account: %s", err.Error()),
			},
		}
	}

	if account == nil {
		return CommandExecutionResult[dto.AuthData]{
			Status: ExecutionStatusFailed,
			Error: CommandErrorDetail{
				Code:    GrantAdminRoleErrorFailedToVerifyAccount,
				Message: fmt.Sprintf("Failed to Verify Account: Account Not Found"),
			},
		}
	}

	account.Roles = []*entities.Role{&entities.Admin}
	_, err = ctx.AccountRepository().Save(account)

	if err != nil {
		return CommandExecutionResult[dto.AuthData]{
			Status: ExecutionStatusFailed,
			Error: CommandErrorDetail{
				Code:    GrantAdminRoleErrorFailedToVerifyAccount,
				Message: fmt.Sprintf("Failed to Update Account: %s", err.Error()),
			},
		}
	}

	return CommandExecutionResult[dto.AuthData]{Status: ExecutionStatusSuccess}
}
