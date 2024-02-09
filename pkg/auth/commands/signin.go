package commands

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/Goldwin/ies-pik-cms/pkg/auth/dto"
	"github.com/Goldwin/ies-pik-cms/pkg/auth/entities"
	. "github.com/Goldwin/ies-pik-cms/pkg/common/commands"
	"github.com/golang-jwt/jwt"
)

type SigninMethod string

const (
	SigninMethodPassword SigninMethod = "password"
	SigninMethodOTP      SigninMethod = "otp"

	SigninErrorFailedToGetOtp                   CommandErrorCode = 20101
	SigninErrorOtpExired                        CommandErrorCode = 20102
	SigninErrorWrongOtp                         CommandErrorCode = 20103
	SigninErrorFailedToCreateToken              CommandErrorCode = 20104
	SigninErrorOTPNotFound                      CommandErrorCode = 20105
	SigninErrorOTPFailedToInvalidateOTP         CommandErrorCode = 20106
	SignInErrorPasswordFailedToGetAccountDetail CommandErrorCode = 20107
	SignInErrorFailedToCreateNewAccount         CommandErrorCode = 20108
)

type SigninCommand struct {
	Email     string
	Password  []byte
	Method    SigninMethod
	SecretKey []byte
}

func (cmd SigninCommand) Execute(ctx CommandContext) CommandExecutionResult[dto.SignInResult] {
	if cmd.Method == SigninMethodPassword {
		return cmd.passwordLogin(ctx)
	}
	return cmd.otpSignIn(ctx)
}

func (cmd SigninCommand) passwordLogin(ctx CommandContext) CommandExecutionResult[dto.SignInResult] {
	passwordRepository := ctx.PasswordRepository()
	password, err := passwordRepository.Get(entities.EmailAddress(cmd.Email))
	if err != nil {
		return CommandExecutionResult[dto.SignInResult]{
			Status: ExecutionStatusFailed,
			Error: CommandErrorDetail{
				Code:    SignInErrorPasswordFailedToGetAccountDetail,
				Message: fmt.Sprintf("Failed to Get Account Detail: %s", err.Error()),
			},
		}
	}
	if password == nil {
		return CommandExecutionResult[dto.SignInResult]{
			Status: ExecutionStatusFailed,
			Error: CommandErrorDetail{
				Code:    SigninErrorOTPNotFound,
				Message: "Wrong Email or Password. Please try again.",
			},
		}
	}

	passwordHash := sha256.Sum256(append(cmd.Password, password.Salt...))

	if !bytes.Equal(password.PasswordHash, passwordHash[:]) {
		return CommandExecutionResult[dto.SignInResult]{
			Status: ExecutionStatusFailed,
			Error: CommandErrorDetail{
				Code:    SigninErrorWrongOtp,
				Message: "Wrong Email or Password. Please try again.",
			},
		}
	}

	tokenString, err := cmd.createToken()

	if err != nil {
		return CommandExecutionResult[dto.SignInResult]{
			Status: ExecutionStatusFailed,
			Error: CommandErrorDetail{
				Code:    SigninErrorFailedToCreateToken,
				Message: fmt.Sprintf("Error When Creating Access Token: %s", err.Error()),
			},
		}
	}

	account, err := ctx.AccountRepository().GetAccount(entities.EmailAddress(cmd.Email))

	if err != nil {
		return CommandExecutionResult[dto.SignInResult]{
			Status: ExecutionStatusFailed,
			Error: CommandErrorDetail{
				Code:    SignInErrorPasswordFailedToGetAccountDetail,
				Message: fmt.Sprintf("Failed to Get Account Detail: %s", err.Error()),
			},
		}
	}

	return CommandExecutionResult[dto.SignInResult]{
		Status: ExecutionStatusSuccess,
		Result: dto.SignInResult{
			AccessToken: tokenString,
			AuthData:    toAuthData(account),
		},
	}
}

func (cmd SigninCommand) otpSignIn(ctx CommandContext) CommandExecutionResult[dto.SignInResult] {
	otpRepository := ctx.OtpRepository()
	accountRepository := ctx.AccountRepository()
	otp, err := otpRepository.GetOtp(entities.EmailAddress(cmd.Email))
	if err != nil {
		return CommandExecutionResult[dto.SignInResult]{
			Status: ExecutionStatusFailed,
			Error: CommandErrorDetail{
				Code:    SigninErrorFailedToGetOtp,
				Message: fmt.Sprintf("Failed to Get OTP: %s", err.Error()),
			},
		}
	}
	if otp == nil {
		return CommandExecutionResult[dto.SignInResult]{
			Status: ExecutionStatusFailed,
			Error: CommandErrorDetail{
				Code:    SigninErrorOTPNotFound,
				Message: "OTP Not Found. Please try again.",
			},
		}
	}

	if otp.ExpiredTime.Before(time.Now()) {
		return CommandExecutionResult[dto.SignInResult]{
			Status: ExecutionStatusFailed,
			Error: CommandErrorDetail{
				Code:    SigninErrorOtpExired,
				Message: "OTP Expired. Please try again.",
			},
		}
	}

	passwordAndSalt := append(cmd.Password, otp.Salt...)
	passwordHash := sha256.Sum256(passwordAndSalt)

	isMatching := bytes.Equal(passwordHash[:], otp.PasswordHash[:])

	if !isMatching {
		return CommandExecutionResult[dto.SignInResult]{
			Status: ExecutionStatusFailed,
			Error: CommandErrorDetail{
				Code:    SigninErrorWrongOtp,
				Message: "Wrong Email Or Wrong Password. Please try again.",
			},
		}
	}

	tokenString, err := cmd.createToken()

	if err != nil {
		return CommandExecutionResult[dto.SignInResult]{
			Status: ExecutionStatusFailed,
			Error: CommandErrorDetail{
				Code:    SigninErrorFailedToGetOtp,
				Message: fmt.Sprintf("Error When Requesting Access Token: %s", err.Error()),
			},
		}
	}

	err = otpRepository.RemoveOtp(*otp)
	if err != nil {
		return CommandExecutionResult[dto.SignInResult]{
			Status: ExecutionStatusFailed,
			Error: CommandErrorDetail{
				Code:    SigninErrorOTPFailedToInvalidateOTP,
				Message: fmt.Sprintf("Failed to invalidate OTP: %s", err.Error()),
			},
		}
	}

	err = otpRepository.RemoveOtp(*otp)

	account, err := accountRepository.GetAccount(entities.EmailAddress(cmd.Email))

	if err != nil {
		return CommandExecutionResult[dto.SignInResult]{
			Status: ExecutionStatusFailed,
			Error: CommandErrorDetail{
				Code:    SigninErrorFailedToGetOtp,
				Message: fmt.Sprintf("Error When Requesting Access Token: %s", err.Error()),
			},
		}
	}

	if account == nil {
		account, err = accountRepository.AddAccount(entities.Account{Email: entities.EmailAddress(cmd.Email)})
		if err != nil {
			return CommandExecutionResult[dto.SignInResult]{
				Status: ExecutionStatusFailed,
				Error: CommandErrorDetail{
					Code:    SignInErrorFailedToCreateNewAccount,
					Message: fmt.Sprintf("Failed to create new account: %s", err.Error()),
				},
			}
		}
	}

	return CommandExecutionResult[dto.SignInResult]{
		Status: ExecutionStatusSuccess,
		Result: dto.SignInResult{
			AccessToken: tokenString,
			AuthData:    toAuthData(account),
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

func toAuthData(account *entities.Account) dto.AuthData {
	scopes := make([]string, 0)

	for _, role := range account.Roles {
		for _, scope := range role.Scopes {
			scopes = append(scopes, string(scope))
		}
	}

	return dto.AuthData{
		ID:         account.Person.ID,
		FirstName:  account.Person.FirstName,
		MiddleName: account.Person.MiddleName,
		LastName:   account.Person.MiddleName,
		Email:      string(account.Email),
		Scopes:     scopes,
	}
}
