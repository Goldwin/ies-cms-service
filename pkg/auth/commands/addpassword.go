package commands

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"math/big"

	"github.com/Goldwin/ies-pik-cms/pkg/auth/dto"
	"github.com/Goldwin/ies-pik-cms/pkg/auth/entities"
	"github.com/Goldwin/ies-pik-cms/pkg/auth/repositories"
	. "github.com/Goldwin/ies-pik-cms/pkg/common/commands"
)

const (
	AddPasswordErrorInvalidInput          AppErrorCode = 20401
	AddPasswordErrorFailedToVerifyAccount AppErrorCode = 20402
)

type SavePasswordCommand struct {
	Input dto.PasswordInput
}

func (cmd SavePasswordCommand) Execute(ctx repositories.CommandContext) AppExecutionResult[dto.PasswordResult] {

	if !bytes.Equal(cmd.Input.Password, cmd.Input.ConfirmPassword) {
		return AppExecutionResult[dto.PasswordResult]{
			Status: ExecutionStatusFailed,
			Error: AppErrorDetail{
				Code:    AddPasswordErrorInvalidInput,
				Message: "Password and Confirm Password should be same",
			},
		}
	}

	account, err := ctx.AccountRepository().GetAccount(entities.EmailAddress(cmd.Input.Email))

	if err != nil {

		return AppExecutionResult[dto.PasswordResult]{
			Status: ExecutionStatusFailed,
			Error: AppErrorDetail{
				Code:    AddPasswordErrorFailedToVerifyAccount,
				Message: fmt.Sprintf("Failed to Verify Account: %s", err.Error()),
			},
		}
	}

	if account == nil {
		_, err = ctx.AccountRepository().AddAccount(entities.Account{
			Email:  entities.EmailAddress(cmd.Input.Email),
			Roles:  []entities.Role{},
			Person: entities.Person{},
		})

		if err != nil {

			return AppExecutionResult[dto.PasswordResult]{
				Status: ExecutionStatusFailed,
				Error: AppErrorDetail{
					Code:    AddPasswordErrorFailedToVerifyAccount,
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
		return AppExecutionResult[dto.PasswordResult]{
			Status: ExecutionStatusFailed,
			Error: AppErrorDetail{
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

	return AppExecutionResult[dto.PasswordResult]{
		Status: ExecutionStatusSuccess,
	}
}
