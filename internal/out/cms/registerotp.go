package cms

import (
	"bytes"
	"fmt"
	"html/template"
	"log"

	"github.com/Goldwin/ies-pik-cms/internal/infra"
	"github.com/Goldwin/ies-pik-cms/pkg/auth/dto"
	"github.com/Goldwin/ies-pik-cms/pkg/common/out"
)

const (
	RegistrationOtpResultSubjectTemplate = "Email Confirmation for %s"
	RegistrationOtpResultMessageTemplate = `    
  <!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN"
  "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
  <html>
    <body>
      <table
        cellpadding="0"
        cellspacing="0"
        style="border-spacing: 0 !important; width: 100%; border-width: 0"
      >
        <tbody>
          <tr>
            <td align="center">
              <table
                cellpadding="0"
                cellspacing="0"
                style="
                  border-spacing: 0 !important;
                  max-width: 600px;
                  width: 100%;
                  border-width: 0;
                "
                bgcolor="#FFFFFF"
              >
                <tbody>
                  <tr>
                    <td
                      align="left"
                      style="word-break: break-word; padding: 64px 32px 32px"
                    >
                      <h1
                        style="
                          color: #333333;
                          font-size: 28px;
                          font-weight: 600;
                          margin: 0 0 1em;
                        "
                      >
                        Your Email has been confirmed.
                      </h1>
                      <p>
                        Here’s the One Time Password (OTP) for the registration                                            

                      <table
                        cellpadding="0"
                        cellspacing="0"
                        style="
                          border-spacing: 0 !important;
                          width: 100%;
                          border-width: 0;
                        "
                      >
                        <tbody>
                          <tr>
                            <td
                              style="border-radius: 3px; padding: 12px"
                              bgcolor="#EAF1FC"
                            >
                              <div style="font-size: 48px" align="center">
                                <strong>{{.OTP}}</strong>
                              </div>
                            </td>
                          </tr>
                        </tbody>
                      </table>
                      <br />
                      <p>
                        If you didn't register, please ignore this email.
                      </p>
                    </td>
                  </tr>
                </tbody>
              </table>
            </td>
          </tr>

          <tr>
            <td align="center">
              <table
                cellpadding="0"
                cellspacing="0"
                style="
                  border-spacing: 0 !important;
                  max-width: 600px;
                  width: 100%;
                  padding-bottom: 24px;
                  padding-top: 24px;
                  border-width: 0;
                "
                bgcolor="#FFFFFF"
              >
                <tbody>
                  <tr>
                    <td
                      colspan="2"
                      style="
                        padding-bottom: 24px;
                        padding-left: 32px;
                        padding-right: 32px;
                        word-break: break-word;
                      "
                    >
                      <table
                        cellpadding="0"
                        cellspacing="0"
                        style="
                          border-spacing: 0 !important;
                          width: 100%;
                          border-width: 0;
                        "
                      >
                        <tbody>
                          <tr>
                            <td
                              style="
                                border-top-width: 1px;
                                border-top-color: #cecece;
                                border-top-style: solid;
                              "
                            ></td>
                          </tr>
                        </tbody>
                      </table>
                    </td>
                  </tr>
                  <tr>
                    <td align="left" style="padding-left: 32px"></td>
                  </tr>
                  <tr>
                    <td
                      align="left"
                      colspan="2"
                      style="word-break: break-word; padding: 24px 32px 48px"
                    >
                      <div
                        style="font-size: 10px; line-height: 1; padding-top: 24px"
                      >
                        You are receiving this communication because you initiated a registration to IES
                        CMS using <a href="mailto:{{.Email}}" target="_blank">{{.Email}}</a>.
                      </div>
                    </td>
                  </tr>
                </tbody>
              </table>
            </td>
          </tr>
        </tbody>
      </table>
    </body>
  </html>
  `
)

type registrationOTPOutputHandler struct {
	emailClient infra.EmailClient
	template    *template.Template
}

// OnError implements out.Output.
func (*registrationOTPOutputHandler) OnError(err out.AppErrorDetail) {
	fmt.Println(err.Error())
}

type OTPData struct {
	Email string
	OTP   string
}

// OnSuccess implements out.Output.
func (o *registrationOTPOutputHandler) OnSuccess(result dto.OtpResult) {
	to := []string{result.Email}
	subject := fmt.Sprintf(RegistrationOtpResultSubjectTemplate, string(result.Email))
	buf := new(bytes.Buffer)
	otpData := OTPData{
		Email: result.Email,
		OTP:   string(result.OTP),
	}
	if err := o.template.Execute(buf, otpData); err != nil {
		log.Default().Printf("Failed to generate email from template: %s", err.Error())
	}
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	message := buf.String()

	err := o.emailClient.Send(to, subject, mime, message)

	if err == nil {
		log.Println("Mail sent!")
	} else {
		log.Println(err)
	}
}

func NewRegistrationOTPOutputHandler(emailClient infra.EmailClient) out.Output[dto.OtpResult] {
	template, err := template.New("registration-otp").Parse(RegistrationOtpResultMessageTemplate)
	if err != nil {
		log.Fatal(err)
	}
	return &registrationOTPOutputHandler{
		emailClient: emailClient,
		template:    template,
	}
}
