package auth

import (
	"fmt"
	"log"

	"github.com/Goldwin/ies-pik-cms/internal/infra"
	"github.com/Goldwin/ies-pik-cms/pkg/auth/dto"
	"github.com/Goldwin/ies-pik-cms/pkg/common/out"
)

const (
	OTPSubjectTemplate = "Your OTP is %s"
	OTPMessageTemplate = "Your OTP is %s"
)

type otpOutputHandler struct {
	emailClient infra.EmailClient
}

// OnError implements out.Output.
func (*otpOutputHandler) OnError(err out.AppErrorDetail) {
	fmt.Println(err.Error())
}

// OnSuccess implements out.Output.
func (o *otpOutputHandler) OnSuccess(result dto.OtpResult) {
	fmt.Println(string(result.OTP))
	to := []string{result.Email}
	subject := fmt.Sprintf(OTPSubjectTemplate, string(result.OTP))
	message := fmt.Sprintf(OTPMessageTemplate, string(result.OTP))

	err := o.emailClient.Send(to, subject, "", message)

	if err == nil {
		log.Println("Mail sent!")
	} else {
		log.Println(err)
	}
}

func newOtpOutputHandler(emailClient infra.EmailClient) out.Output[dto.OtpResult] {
	return &otpOutputHandler{
		emailClient: emailClient,
	}
}
