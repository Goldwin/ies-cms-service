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

type SavePasswordCommand struct {
	Email    entities.EmailAddress
	Password []byte
}

func (cmd SavePasswordCommand) Execute(ctx CommandContext) CommandExecutionResult[dto.PasswordResult] {

	account, err := ctx.AccountRepository().Get(string(cmd.Email))

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
		_, err = ctx.AccountRepository().Save(&entities.Account{
			Email: entities.EmailAddress(cmd.Email),
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
		EmailAddress: entities.EmailAddress(cmd.Email),
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
	passwordAndSalt := append(cmd.Password, password.Salt...)
	passwordHash := sha256.Sum256(passwordAndSalt)
	password.PasswordHash = passwordHash[:]

	err = ctx.PasswordRepository().Save(password)
	if err != nil {
		return CommandExecutionResult[dto.PasswordResult]{
			Status: ExecutionStatusFailed,
			Error: CommandErrorDetail{
				Code:    SavePasswordErrorFailedToVerifyAccount,
				Message: fmt.Sprintf("Failed to Save Password: %s", err.Error()),
			},
		}
	}

	return CommandExecutionResult[dto.PasswordResult]{
		Status: ExecutionStatusSuccess,
		Result: dto.PasswordResult{
			Email: string(password.EmailAddress),
		},
	}
}
