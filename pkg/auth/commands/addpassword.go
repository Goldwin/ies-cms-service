package commands

import (
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
	AddPasswordErrorInvalidInput AppErrorCode = 20401
)

type AddPassword struct {
	input dto.PasswordInput
}

func (cmd AddPassword) Execute(ctx repositories.CommandContext) AppExecutionResult[dto.PasswordResult] {

	if cmd.input.Password != cmd.input.ConfirmPassword {
		return AppExecutionResult[dto.PasswordResult]{
			Status: ExecutionStatusFailed,
			Error: AppErrorDetail{
				Code:    AddPasswordErrorInvalidInput,
				Message: "Password and Confirm Password should be same",
			},
		}
	}
	password := entities.PasswordDetail{
		EmailAddress: entities.EmailAddress(cmd.input.Email),
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
	passwordAndSalt := append([]byte(cmd.input.Email), password.Salt...)
	passwordHash := sha256.Sum256(passwordAndSalt)
	password.PasswordHash = passwordHash[:]

	err = ctx.PasswordRepository().Save(password)

	return AppExecutionResult[dto.PasswordResult]{
		Status: ExecutionStatusSuccess,
	}
}
