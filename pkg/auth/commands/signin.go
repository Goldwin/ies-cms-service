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
	SecretKey []byte
}

func (cmd SigninCommand) Execute(ctx CommandContext) CommandExecutionResult[dto.SignInResult] {
	return cmd.passwordLogin(ctx)
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

	account, err := ctx.AccountRepository().Get(cmd.Email)

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
		Email:  string(account.Email),
		Scopes: scopes,
	}
}
