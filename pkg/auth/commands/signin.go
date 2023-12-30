package commands

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/Goldwin/ies-pik-cms/pkg/auth/dto"
	"github.com/Goldwin/ies-pik-cms/pkg/auth/entities"
	"github.com/Goldwin/ies-pik-cms/pkg/auth/repositories"
	. "github.com/Goldwin/ies-pik-cms/pkg/common/commands"
	"github.com/golang-jwt/jwt"
)

type SigninMethod string

const (
	SigninMethodPassword SigninMethod = "password"
	SigninMethodOTP      SigninMethod = "otp"

	SigninErrorFailedToGetOtp            AppErrorCode = 20101
	SigninErrorOtpExired                 AppErrorCode = 20102
	SigninErrorWrongOtp                  AppErrorCode = 20103
	SigninErrorFailedToCreateToken       AppErrorCode = 20104
	SigninErrorOTPNotFound               AppErrorCode = 20105
	SigninErrorOTPFailedToInvalidateOTP  AppErrorCode = 20106
	SignInErrorPasswordLoginNotSupported AppErrorCode = 20107
	SignInErrorFailedToCreateNewAccount  AppErrorCode = 20108
)

type SigninCommand struct {
	Email     string
	Password  []byte
	Method    SigninMethod
	SecretKey []byte
}

func (cmd SigninCommand) Execute(ctx repositories.CommandContext) AppExecutionResult[dto.SignInResult] {
	if cmd.Method == SigninMethodPassword {
		return AppExecutionResult[dto.SignInResult]{
			Status: ExecutionStatusFailed,
			Error: AppErrorDetail{
				Code:    SignInErrorPasswordLoginNotSupported,
				Message: "Password Login Not Supported",
			},
		}
	}
	return cmd.otpSignIn(ctx.OtpRepository(), ctx.AccountRepository())
}

func (cmd SigninCommand) otpSignIn(otpRepository repositories.OtpRepository, accountRepository repositories.AccountRepository) AppExecutionResult[dto.SignInResult] {
	otp, err := otpRepository.GetOtp(entities.EmailAddress(cmd.Email))
	if err != nil {
		return AppExecutionResult[dto.SignInResult]{
			Status: ExecutionStatusFailed,
			Error: AppErrorDetail{
				Code:    SigninErrorFailedToGetOtp,
				Message: fmt.Sprintf("Failed to Get OTP: %s", err.Error()),
			},
		}
	}
	if otp == nil {
		return AppExecutionResult[dto.SignInResult]{
			Status: ExecutionStatusFailed,
			Error: AppErrorDetail{
				Code:    SigninErrorOTPNotFound,
				Message: "OTP Not Found. Please try again.",
			},
		}
	}

	if otp.ExpiredTime.Before(time.Now()) {
		return AppExecutionResult[dto.SignInResult]{
			Status: ExecutionStatusFailed,
			Error: AppErrorDetail{
				Code:    SigninErrorOtpExired,
				Message: "OTP Expired. Please try again.",
			},
		}
	}

	passwordAndSalt := append(cmd.Password, otp.Salt...)
	passwordHash := sha256.Sum256(passwordAndSalt)

	isMatching := bytes.Equal(passwordHash[:], otp.PasswordHash[:])

	if !isMatching {
		return AppExecutionResult[dto.SignInResult]{
			Status: ExecutionStatusFailed,
			Error: AppErrorDetail{
				Code:    SigninErrorWrongOtp,
				Message: "Wrong Email Or Wrong Password. Please try again.",
			},
		}
	}

	tokenString, err := cmd.createToken()

	if err != nil {
		return AppExecutionResult[dto.SignInResult]{
			Status: ExecutionStatusFailed,
			Error: AppErrorDetail{
				Code:    SigninErrorFailedToGetOtp,
				Message: fmt.Sprintf("Error When Requesting Access Token: %s", err.Error()),
			},
		}
	}

	err = otpRepository.RemoveOtp(*otp)
	if err != nil {
		return AppExecutionResult[dto.SignInResult]{
			Status: ExecutionStatusFailed,
			Error: AppErrorDetail{
				Code:    SigninErrorOTPFailedToInvalidateOTP,
				Message: fmt.Sprintf("Failed to invalidate OTP: %s", err.Error()),
			},
		}
	}

	err = otpRepository.RemoveOtp(*otp)

	account, err := accountRepository.GetAccount(entities.EmailAddress(cmd.Email))

	if err != nil {
		return AppExecutionResult[dto.SignInResult]{
			Status: ExecutionStatusFailed,
			Error: AppErrorDetail{
				Code:    SigninErrorFailedToGetOtp,
				Message: fmt.Sprintf("Error When Requesting Access Token: %s", err.Error()),
			},
		}
	}

	if account == nil {
		account, err = accountRepository.AddAccount(entities.Account{Email: entities.EmailAddress(cmd.Email)})
		if err != nil {
			return AppExecutionResult[dto.SignInResult]{
				Status: ExecutionStatusFailed,
				Error: AppErrorDetail{
					Code:    SignInErrorFailedToCreateNewAccount,
					Message: fmt.Sprintf("Failed to create new account: %s", err.Error()),
				},
			}
		}
	}

	return AppExecutionResult[dto.SignInResult]{
		Status: ExecutionStatusSuccess,
		Result: dto.SignInResult{
			AccessToken: tokenString,
			AuthData: dto.AuthData{
				ID:         account.Person.ID,
				FirstName:  account.Person.FirstName,
				MiddleName: account.Person.MiddleName,
				LastName:   account.Person.MiddleName,
				Email:      cmd.Email,
				Scopes:     []string{},
			},
		},
	}
}

func (cmd SigninCommand) createToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"email": cmd.Email,
			"exp":   time.Now().Add(time.Hour * 24 * 14).Unix(),
		})

	tokenString, err := token.SignedString(cmd.SecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
