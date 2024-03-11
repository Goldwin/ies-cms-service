package commands

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"math/big"

	"github.com/Goldwin/ies-pik-cms/pkg/auth/dto"
	"github.com/Goldwin/ies-pik-cms/pkg/auth/entities"
	. "github.com/Goldwin/ies-pik-cms/pkg/common/commands"
)

const (
	SavePasswordErrorInvalidInput          CommandErrorCode = 20401
	SavePasswordErrorFailedToVerifyAccount CommandErrorCode = 20402
)

type ResetPasswordCommand struct {
	Input dto.PasswordInput
}

func (cmd ResetPasswordCommand) Execute(ctx CommandContext) CommandExecutionResult[dto.PasswordResult] {

	token, err := ctx.PasswordRepository().GetResetToken(entities.EmailAddress(cmd.Input.Email))

	if err != nil {
		return CommandExecutionResult[dto.PasswordResult]{
			Status: ExecutionStatusFailed,
			Error: CommandErrorDetail{
				Code:    SavePasswordErrorInvalidInput,
				Message: fmt.Sprintf("Failed to fetch Reset Token: %s", err.Error()),
			},
		}
	}

	if token != cmd.Input.Token {
		return CommandExecutionResult[dto.PasswordResult]{
			Status: ExecutionStatusFailed,
			Error: CommandErrorDetail{
				Code:    SavePasswordErrorInvalidInput,
				Message: fmt.Sprintf("Reset Token Mismatched"),
			},
		}
	}

	account, err := ctx.AccountRepository().GetAccount(entities.EmailAddress(cmd.Input.Email))

	if err != nil {

		return CommandExecutionResult[dto.PasswordResult]{
			Status: ExecutionStatusFailed,
			Error: CommandErrorDetail{
				Code:    SavePasswordErrorFailedToVerifyAccount,
				Message: fmt.Sprintf("Failed to Verify Account: %s", err.Error()),
			},
		}
	}

	if account == nil {
		_, err = ctx.AccountRepository().AddAccount(entities.Account{
			Email: entities.EmailAddress(cmd.Input.Email),
			Roles: []entities.Role{},
		})

		if err != nil {

			return CommandExecutionResult[dto.PasswordResult]{
				Status: ExecutionStatusFailed,
				Error: CommandErrorDetail{
					Code:    SavePasswordErrorFailedToVerifyAccount,
					Message: fmt.Sprintf("Failed to Verify Account: %s", err.Error()),
				},
			}
		}
	}

	password := entities.PasswordDetail{
		EmailAddress: entities.EmailAddress(cmd.Input.Email),
	}

	salt, err := rand.Int(rand.Reader, big.NewInt(999999))
	if err != nil {
		return CommandExecutionResult[dto.PasswordResult]{
			Status: ExecutionStatusFailed,
			Error: CommandErrorDetail{
				Code:    GenerateOtpErrorFailedToGenOtp,
				Message: fmt.Sprintf("Failed to Save Password: %s", err.Error()),
			},
		}
	}

	password.Salt = salt.Bytes()
	passwordAndSalt := append([]byte(cmd.Input.Password), password.Salt...)
	passwordHash := sha256.Sum256(passwordAndSalt)
	password.PasswordHash = passwordHash[:]

	err = ctx.PasswordRepository().Save(password)

	return CommandExecutionResult[dto.PasswordResult]{
		Status: ExecutionStatusSuccess,
	}
}
