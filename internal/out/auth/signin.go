package auth

import (
	"github.com/Goldwin/ies-pik-cms/pkg/auth/dto"
	"github.com/Goldwin/ies-pik-cms/pkg/common/commands"
	"github.com/Goldwin/ies-pik-cms/pkg/common/out"
)

type signInOutputHandler struct {
}

// OnError implements out.Output.
func (*signInOutputHandler) OnError(err commands.AppErrorDetail) {
}

// OnSuccess implements out.Output.
func (o *signInOutputHandler) OnSuccess(result dto.SignInResult) {
}

func newSampleSignInOutputHandler() out.Output[dto.SignInResult] {
	return &signInOutputHandler{}
}
