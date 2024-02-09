package auth

import (
	"fmt"
	"log"
	"net/smtp"
	"os"

	"github.com/Goldwin/ies-pik-cms/pkg/auth/dto"
	"github.com/Goldwin/ies-pik-cms/pkg/common/out"
)

type otpOutputHandler struct {
	SenderEmail     string
	SenderName      string
	SenderPassword  []byte
	SMTPHost        string
	SMTPPort        int
	SubjectTemplate string
	MessageTemplate string
}

// OnError implements out.Output.
func (*otpOutputHandler) OnError(err out.AppErrorDetail) {
	fmt.Println(err.Error())
}

// OnSuccess implements out.Output.
func (o *otpOutputHandler) OnSuccess(result dto.OtpResult) {
	fmt.Println(string(result.OTP))
	to := []string{result.Email}
	subject := fmt.Sprintf(o.SubjectTemplate, string(result.OTP))
	message := fmt.Sprintf(o.MessageTemplate, string(result.OTP))
	hostAndPort := fmt.Sprintf("%s:%d", o.SMTPHost, o.SMTPPort)

	auth := smtp.PlainAuth("", o.SenderEmail, "aanxgrufighqhptv", o.SMTPHost)

	msg := fmt.Sprintf("Subject: %s\n\n%s", subject, message)
	err := smtp.SendMail(hostAndPort, auth, o.SenderEmail, to, []byte(msg))

	if err != nil {
		log.Fatal(err.Error())
	}

	log.Println("Mail sent!")
}

func newOtpOutputHandler() out.Output[dto.OtpResult] {
	//TODO fix this.
	return &otpOutputHandler{
		SenderEmail:     os.Getenv("EMAIL_SENDER_ADDRESS"),
		SenderName:      fmt.Sprintf("no-reply <%s>", os.Getenv("EMAIL_SENDER_ADDRESS")),
		SenderPassword:  []byte(os.Getenv("EMAIL_SENDER_PASSWORD")),
		SMTPHost:        "smtp.gmail.com",
		SMTPPort:        587,
		SubjectTemplate: "%v is your OTP for IES Pik App",
		MessageTemplate: "OTP: %v",
	}
}
