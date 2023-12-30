package auth

import (
	"github.com/Goldwin/ies-pik-cms/internal/bus"
	"github.com/Goldwin/ies-pik-cms/internal/out/common"
	"github.com/Goldwin/ies-pik-cms/pkg/auth/dto"
	"github.com/Goldwin/ies-pik-cms/pkg/common/out"
)

type AuthOutputComponent interface {
	OTPOutput() out.Output[dto.OtpResult]
	SignInOutput() out.Output[dto.SignInResult]
	AuthOutput() out.Output[dto.AuthData]
	RegistrationOutput() out.Output[dto.AuthData]
}

type authOutputComponentImpl struct {
	otpOutput             out.Output[dto.OtpResult]
	signInOutputHandler   out.Output[dto.SignInResult]
	authOutputHandler     out.Output[dto.AuthData]
	registerOutputHandler out.Output[dto.AuthData]
}

// RegistrationOutput implements AuthOutputComponent.
func (a *authOutputComponentImpl) RegistrationOutput() out.Output[dto.AuthData] {
	return a.registerOutputHandler
}

// AuthOutput implements AuthOutputComponent.
func (*authOutputComponentImpl) AuthOutput() out.Output[dto.AuthData] {
	return &common.NoopOutput[dto.AuthData]{}
}

// SignInOutput implements AuthOutputComponent.
func (a *authOutputComponentImpl) SignInOutput() out.Output[dto.SignInResult] {
	return a.signInOutputHandler
}

// OTPOutput implements AuthOutputComponent.
func (a *authOutputComponentImpl) OTPOutput() out.Output[dto.OtpResult] {
	return a.otpOutput
}

func NewAuthOutputComponent(eventBus bus.EventBusComponent) AuthOutputComponent {
	return &authOutputComponentImpl{
		otpOutput:             newOtpOutputHandler(),
		signInOutputHandler:   newSampleSignInOutputHandler(),
		registerOutputHandler: newRegisterOutputHandler(eventBus),
	}
}
