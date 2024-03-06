package auth

import (
	"context"
	"log"
	"os"

	"github.com/Goldwin/ies-pik-cms/pkg/auth/commands"
	"github.com/Goldwin/ies-pik-cms/pkg/auth/dto"
	"github.com/Goldwin/ies-pik-cms/pkg/common"
	. "github.com/Goldwin/ies-pik-cms/pkg/common/commands"
	"github.com/Goldwin/ies-pik-cms/pkg/common/out"
	"github.com/Goldwin/ies-pik-cms/pkg/common/worker"
)

type AuthDataLayerComponent interface {
	CommandWorker() worker.UnitOfWork[commands.CommandContext]
}

type AuthComponent interface {
	SignIn(ctx context.Context, input dto.SignInInput, output out.Output[dto.SignInResult])
	GenerateOtp(ctx context.Context, input dto.OtpInput, output out.Output[dto.OtpResult])
	CompleteRegistration(ctx context.Context, input dto.CompleteRegistrationInput, output out.Output[dto.AuthData])
	Auth(ctx context.Context, input dto.AuthInput, output out.Output[dto.AuthData])
	SavePassword(ctx context.Context, input dto.PasswordInput, output out.Output[dto.PasswordResult])
	GenerateResetToken(ctx context.Context, email string, output out.Output[string])
	common.Component
}

type authComponentImpl struct {
	worker    worker.UnitOfWork[commands.CommandContext]
	secretKey []byte
}

// GenerateResetToken implements AuthComponent.
func (a *authComponentImpl) GenerateResetToken(ctx context.Context, email string, output out.Output[string]) {
	var res CommandExecutionResult[string]
	a.worker.Execute(context.Background(), func(ctx commands.CommandContext) error {
		res = commands.GenerateResetTokenCommand{
			Email: email,
		}.Execute(ctx)
		if res.Status != ExecutionStatusSuccess {
			return res.Error
		}
		return nil
	})

	if res.Status == ExecutionStatusSuccess {
		output.OnSuccess(res.Result)
	} else {
		output.OnError(out.ConvertCommandErrorDetail(res.Error))
	}
}

// SavePassword implements AuthComponent.
func (a *authComponentImpl) SavePassword(ctx context.Context, input dto.PasswordInput, output out.Output[dto.PasswordResult]) {
	var res CommandExecutionResult[dto.PasswordResult]
	a.worker.Execute(context.Background(), func(ctx commands.CommandContext) error {
		res = commands.SavePasswordCommand{
			Input: input,
		}.Execute(ctx)
		if res.Status != ExecutionStatusSuccess {
			return res.Error
		}
		return nil
	})

	if res.Status == ExecutionStatusSuccess {
		output.OnSuccess(res.Result)
	} else {
		output.OnError(out.ConvertCommandErrorDetail(res.Error))
	}
}

// Start implements AuthComponent.
func (a *authComponentImpl) Start() {
	a.worker.Execute(context.Background(), func(ctx commands.CommandContext) error {
		res := commands.SavePasswordCommand{
			Input: dto.PasswordInput{
				Email:           os.Getenv("ADMIN_EMAIL"),
				Password:        []byte(os.Getenv("ADMIN_PASSWORD")),
				ConfirmPassword: []byte(os.Getenv("ADMIN_PASSWORD")),
			},
		}.Execute(ctx)
		if res.Status != ExecutionStatusSuccess {
			log.Fatal(res.Error)
		}
		return nil
	})
	a.worker.Execute(context.Background(), func(ctx commands.CommandContext) error {
		res := commands.GrantAdminRoleCommand{
			Email: os.Getenv("ADMIN_EMAIL"),
		}.Execute(ctx)

		if res.Status != ExecutionStatusSuccess {
			log.Fatal(res.Error)
		}
		return nil
	})
}

// Stop implements AuthComponent.
func (a *authComponentImpl) Stop() {
	log.Default().Printf("Auth Component Stopped!")
}

// Auth implements AuthComponent.
func (a *authComponentImpl) Auth(ctx context.Context, input dto.AuthInput, output out.Output[dto.AuthData]) {
	var result CommandExecutionResult[dto.AuthData]
	_ = a.worker.Execute(ctx, func(ctx commands.CommandContext) error {
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
		output.OnError(out.ConvertCommandErrorDetail(result.Error))
	}

}

// CompleteRegistration implements AuthComponent.
func (a *authComponentImpl) CompleteRegistration(ctx context.Context, input dto.CompleteRegistrationInput, output out.Output[dto.AuthData]) {
	var result CommandExecutionResult[dto.AuthData]
	_ = a.worker.Execute(ctx, func(ctx commands.CommandContext) error {
		result = commands.CompleteRegistrationCommand{
			Input: input,
		}.Execute(ctx)
		if result.Status != ExecutionStatusSuccess {
			return result.Error
		}
		return nil
	})
	if result.Status == ExecutionStatusSuccess {
		output.OnSuccess(result.Result)
	} else {
		output.OnError(out.ConvertCommandErrorDetail(result.Error))
	}
}

// GenerateOtp implements AuthComponent.
func (a *authComponentImpl) GenerateOtp(ctx context.Context, input dto.OtpInput, output out.Output[dto.OtpResult]) {
	var result CommandExecutionResult[dto.OtpResult]
	_ = a.worker.Execute(ctx, func(ctx commands.CommandContext) error {
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
		output.OnError(out.ConvertCommandErrorDetail(result.Error))
	}
}

// SignIn implements AuthComponent.
func (a *authComponentImpl) SignIn(ctx context.Context, input dto.SignInInput, output out.Output[dto.SignInResult]) {
	var result CommandExecutionResult[dto.SignInResult]
	_ = a.worker.Execute(ctx, func(ctx commands.CommandContext) error {
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
		output.OnError(out.ConvertCommandErrorDetail(result.Error))
	}
}

func NewAuthComponent(component AuthDataLayerComponent, secretKey []byte) AuthComponent {
	return &authComponentImpl{
		worker:    component.CommandWorker(),
		secretKey: secretKey,
	}
}
