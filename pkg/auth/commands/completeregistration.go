package commands

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/Goldwin/ies-pik-cms/pkg/auth/dto"
	"github.com/Goldwin/ies-pik-cms/pkg/auth/entities"
	. "github.com/Goldwin/ies-pik-cms/pkg/common/commands"
)

const (
	CompleteRegistrationErrorAlreadyRegistered      CommandErrorCode = 20201
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
	accountResult := cmd.verifyOTPAndCreateAccount(ctx)
	if accountResult.Status != ExecutionStatusSuccess {
		return CommandExecutionResult[dto.AuthData]{
			Status: ExecutionStatusFailed,
			Error:  accountResult.Error,
		}
	}

	account, err := ctx.AccountRepository().Get(cmd.Input.Email)
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
		account = &entities.Account{
			Email: entities.EmailAddress(cmd.Input.Email),
		}
	}

	if len(account.Roles) > 0 {
		return CommandExecutionResult[dto.AuthData]{
			Status: ExecutionStatusFailed,
			Error: CommandErrorDetail{
				Code:    CompleteRegistrationErrorAlreadyRegistered,
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

	account.Roles = []*entities.Role{
		&entities.ChurchMember,
	}

	account, err = ctx.AccountRepository().Save(account)

	if err != nil {
		return CommandExecutionResult[dto.AuthData]{
			Status: ExecutionStatusFailed,
			Error: CommandErrorDetail{
				Code:    CompleteRegistrationErrorUpdateFailure,
				Message: fmt.Sprintf("Failed to Update Account: %s", err.Error()),
			},
		}
	}

	result := SavePasswordCommand{
		Email:    entities.EmailAddress(cmd.Input.Email),
		Password: []byte(cmd.Input.Password),
	}.Execute(ctx)

	if result.Status != ExecutionStatusSuccess {
		return CommandExecutionResult[dto.AuthData]{
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

	return CommandExecutionResult[dto.AuthData]{
		Status: ExecutionStatusSuccess,
		Result: dto.AuthData{
			Email:  string(account.Email),
			Scopes: scopes,
		},
	}
}

func (cmd CompleteRegistrationCommand) verifyOTPAndCreateAccount(ctx CommandContext) CommandExecutionResult[*entities.Account] {
	otpRepository := ctx.OtpRepository()
	accountRepository := ctx.AccountRepository()
	otp, err := otpRepository.Get(cmd.Input.Email)
	if err != nil {
		return CommandExecutionResult[*entities.Account]{
			Status: ExecutionStatusFailed,
			Error: CommandErrorDetail{
				Code:    SigninErrorFailedToGetOtp,
				Message: fmt.Sprintf("Failed to Get OTP: %s", err.Error()),
			},
		}
	}
	if otp == nil {
		return CommandExecutionResult[*entities.Account]{
			Status: ExecutionStatusFailed,
			Error: CommandErrorDetail{
				Code:    SigninErrorOTPNotFound,
				Message: "OTP Not Found. Please try again.",
			},
		}
	}

	if otp.ExpiresAt.Before(time.Now()) {
		return CommandExecutionResult[*entities.Account]{
			Status: ExecutionStatusFailed,
			Error: CommandErrorDetail{
				Code:    SigninErrorOtpExired,
				Message: "OTP Expired. Please try again.",
			},
		}
	}

	passwordAndSalt := append([]byte(cmd.Input.OTP), otp.Salt...)
	passwordHash := sha256.Sum256(passwordAndSalt)

	isMatching := bytes.Equal(passwordHash[:], otp.PasswordHash[:])

	if !isMatching {
		return CommandExecutionResult[*entities.Account]{
			Status: ExecutionStatusFailed,
			Error: CommandErrorDetail{
				Code:    SigninErrorWrongOtp,
				Message: "Wrong OTP. Please try again.",
			},
		}
	}

	err = otpRepository.Delete(otp)
	if err != nil {
		return CommandExecutionResult[*entities.Account]{
			Status: ExecutionStatusFailed,
			Error: CommandErrorDetail{
				Code:    SigninErrorOTPFailedToInvalidateOTP,
				Message: fmt.Sprintf("Failed to invalidate OTP: %s", err.Error()),
			},
		}
	}

	err = otpRepository.Delete(otp)

	account, err := accountRepository.Get(cmd.Input.Email)

	if err != nil {
		return CommandExecutionResult[*entities.Account]{
			Status: ExecutionStatusFailed,
			Error: CommandErrorDetail{
				Code:    SigninErrorFailedToGetOtp,
				Message: fmt.Sprintf("Error When Requesting Access Token: %s", err.Error()),
			},
		}
	}

	if account == nil {
		account, err = accountRepository.Save(&entities.Account{Email: entities.EmailAddress(cmd.Input.Email)})
		if err != nil {
			return CommandExecutionResult[*entities.Account]{
				Status: ExecutionStatusFailed,
				Error: CommandErrorDetail{
					Code:    SignInErrorFailedToCreateNewAccount,
					Message: fmt.Sprintf("Failed to create new account: %s", err.Error()),
				},
			}
		}
	}

	return CommandExecutionResult[*entities.Account]{
		Status: ExecutionStatusSuccess,
		Result: account,
	}
}
