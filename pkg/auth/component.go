package auth

import (
	"context"

	"github.com/Goldwin/ies-pik-cms/pkg/auth/commands"
	"github.com/Goldwin/ies-pik-cms/pkg/auth/dto"
	"github.com/Goldwin/ies-pik-cms/pkg/auth/repositories"
	. "github.com/Goldwin/ies-pik-cms/pkg/common/commands"
	"github.com/Goldwin/ies-pik-cms/pkg/common/out"
	"github.com/Goldwin/ies-pik-cms/pkg/common/worker"
)

type AuthDataLayerComponent interface {
	CommandWorker() worker.UnitOfWork[repositories.CommandContext]
}

type AuthComponent interface {
	SignIn(ctx context.Context, input dto.SignInInput, output out.Output[dto.SignInResult])
	GenerateOtp(ctx context.Context, input dto.OtpInput, output out.Output[dto.OtpResult])
	CompleteRegistration(ctx context.Context, input dto.CompleteRegistrationInput, output out.Output[dto.AuthData])
	Auth(ctx context.Context, input dto.AuthInput, output out.Output[dto.AuthData])
}

type authComponentImpl struct {
	worker    worker.UnitOfWork[repositories.CommandContext]
	secretKey []byte
}

// Auth implements AuthComponent.
func (a *authComponentImpl) Auth(ctx context.Context, input dto.AuthInput, output out.Output[dto.AuthData]) {
	var result AppExecutionResult[dto.AuthData]
	_ = a.worker.Execute(ctx, func(ctx repositories.CommandContext) error {
		result = commands.AuthCommand{
			Token:     input.Token,
			SecretKey: a.secretKey,
		}.Execute(ctx)
		if result.Status != ExecutionStatusSuccess {
			return result.Error
		}
		return nil
	})
	if result.Status == ExecutionStatusSuccess {
		output.OnSuccess(result.Result)
	} else {
		output.OnError(result.Error)
	}
}

// CompleteRegistration implements AuthComponent.
func (a *authComponentImpl) CompleteRegistration(ctx context.Context, input dto.CompleteRegistrationInput, output out.Output[dto.AuthData]) {
	var result AppExecutionResult[dto.AuthData]
	_ = a.worker.Execute(ctx, func(ctx repositories.CommandContext) error {
		result = commands.CompleteRegistrationCommand{
			FirstName:  input.FirstName,
			MiddleName: input.MiddleName,
			LastName:   input.LastName,
			Email:      input.Email,
		}.Execute(ctx)
		if result.Status != ExecutionStatusSuccess {
			return result.Error
		}
		return nil
	})
	if result.Status == ExecutionStatusSuccess {
		output.OnSuccess(result.Result)
	} else {
		output.OnError(result.Error)
	}

}

// GenerateOtp implements AuthComponent.
func (a *authComponentImpl) GenerateOtp(ctx context.Context, input dto.OtpInput, output out.Output[dto.OtpResult]) {
	var result AppExecutionResult[dto.OtpResult]
	_ = a.worker.Execute(ctx, func(ctx repositories.CommandContext) error {
		result = commands.GenerateOtpCommand{
			Email: input.Email,
		}.Execute(ctx)
		if result.Status != ExecutionStatusSuccess {
			return result.Error
		}
		return nil
	})

	if result.Status == ExecutionStatusSuccess {
		output.OnSuccess(result.Result)
	} else {
		output.OnError(result.Error)
	}
}

// SignIn implements AuthComponent.
func (a *authComponentImpl) SignIn(ctx context.Context, input dto.SignInInput, output out.Output[dto.SignInResult]) {
	var result AppExecutionResult[dto.SignInResult]
	_ = a.worker.Execute(ctx, func(ctx repositories.CommandContext) error {
		result = commands.SigninCommand{
			Email:     input.Email,
			Password:  []byte(input.Password),
			Method:    commands.SigninMethod(input.Method),
			SecretKey: a.secretKey,
		}.Execute(ctx)
		if result.Status != ExecutionStatusSuccess {
			return result.Error
		}
		return nil
	})
	if result.Status == ExecutionStatusSuccess {
		output.OnSuccess(result.Result)
	} else {
		output.OnError(result.Error)
	}
}

func NewAuthComponent(component AuthDataLayerComponent, secretKey []byte) AuthComponent {
	return &authComponentImpl{
		worker:    component.CommandWorker(),
		secretKey: secretKey,
	}
}
