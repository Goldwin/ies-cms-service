package commands

import (
	"fmt"

	"github.com/Goldwin/ies-pik-cms/pkg/auth/dto"
	"github.com/Goldwin/ies-pik-cms/pkg/auth/entities"
	. "github.com/Goldwin/ies-pik-cms/pkg/common/commands"
	"github.com/golang-jwt/jwt"
)

const (
	AuthErrorInvalidToken            CommandErrorCode = 20301
	AuthErrorFailedToRetrieveAccount CommandErrorCode = 20302
	AuthErrorAccountDoesNotExist     CommandErrorCode = 20303
	AuthErrorOtpExists               CommandErrorCode = 20304
)

type AuthCommand struct {
	Token     string
	SecretKey []byte
}

func (cmd AuthCommand) Execute(ctx CommandContext) CommandExecutionResult[dto.AuthData] {

	claims, err := cmd.extractClaims()

	if err != nil {
		return CommandExecutionResult[dto.AuthData]{
			Status: ExecutionStatusFailed,
			Error: CommandErrorDetail{
				Code:    AuthErrorInvalidToken,
				Message: fmt.Sprintf("Invalid Token: %s", err.Error()),
			},
		}
	}

	emailStr, ok := claims["email"].(string)
	if !ok {
		return CommandExecutionResult[dto.AuthData]{
			Status: ExecutionStatusFailed,
			Error: CommandErrorDetail{
				Code:    AuthErrorInvalidToken,
				Message: fmt.Sprintf("Invalid Token: Malformed Token"),
			},
		}
	}
	email := entities.EmailAddress(emailStr)

	if !email.IsValid() {
		return CommandExecutionResult[dto.AuthData]{
			Status: ExecutionStatusFailed,
			Error: CommandErrorDetail{
				Code:    AuthErrorInvalidToken,
				Message: "Invalid Token: Invalid Email",
			},
		}
	}

	account, err := ctx.AccountRepository().Get(string(email))

	if err != nil {
		return CommandExecutionResult[dto.AuthData]{
			Status: ExecutionStatusFailed,
			Error: CommandErrorDetail{
				Code:    AuthErrorFailedToRetrieveAccount,
				Message: fmt.Sprintf("Invalid Token: %s", err.Error()),
			},
		}
	}

	if account == nil {
		return CommandExecutionResult[dto.AuthData]{
			Status: ExecutionStatusFailed,
			Error: CommandErrorDetail{
				Code:    AuthErrorAccountDoesNotExist,
				Message: "Account Does not exists",
			},
		}
	}

	scopeMap := make(map[entities.Scope]bool, 0)
	scopes := make([]string, 0)

	for _, role := range account.Roles {
		for _, scope := range role.Scopes {
			scopeMap[scope] = true
		}
	}

	for scope := range scopeMap {
		scopes = append(scopes, string(scope))
	}

	return CommandExecutionResult[dto.AuthData]{
		Status: ExecutionStatusSuccess,
		Error:  CommandErrorDetail{},
		Result: dto.AuthData{
			Email:  string(email),
			Scopes: scopes,
		},
	}
}

func (cmd AuthCommand) extractClaims() (jwt.MapClaims, error) {
	token, err := jwt.Parse(cmd.Token, func(token *jwt.Token) (interface{}, error) {
		return cmd.SecretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid token")
}
