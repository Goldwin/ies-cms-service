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
	PasswordResetSubjectTemplate = "Password Reset Request %s"
	PasswordResetMessageTemplate = `    
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
                        A password reset has been initiated.
                      </h1>
                      <p>
                        Here’s a reset code for
                        <a href="mailto:{{.Email}}" target="_blank">{{.Email}}</a>.
                        Enter it
                        <a
                          href="https://iespik.brightfellow.net/password/reset/email_code?login={{.Email}}"
                          style="color: #529ff8; text-decoration: none"
                          target="_blank"
                          data-saferedirecturl=""
                          >on this page</a
                        >
                        to reset your password.
                      </p>

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
                                <strong>{{.Token}}</strong>
                              </div>
                            </td>
                          </tr>
                        </tbody>
                      </table>
                      <br />
                      <p>
                        If you didn't initiate a password reset, just ignore this
                        email.
                      </p>

                      <p>Thanks for using Planning Center!</p>
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
                        You are receiving this communication because you have a IES
                        CMS account and a password reset for
                        <a href="mailto:{{.Email}}" target="_blank">{{.Email}}</a>
                        was requested
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

type passwordResetOutputHandler struct {
	emailClient infra.EmailClient
	template    *template.Template
}

// OnError implements out.Output.
func (*passwordResetOutputHandler) OnError(err out.AppErrorDetail) {
	fmt.Println(err.Error())
}

// OnSuccess implements out.Output.
func (o *passwordResetOutputHandler) OnSuccess(result dto.PasswordResetTokenResult) {
	to := []string{result.Email}
	subject := fmt.Sprintf(PasswordResetSubjectTemplate, string(result.Token))
	buf := new(bytes.Buffer)
	if err := o.template.Execute(buf, result); err != nil {
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

func NewPasswordResetOutputHandler(emailClient infra.EmailClient) out.Output[dto.PasswordResetTokenResult] {
	template, err := template.New("password-reset").Parse(PasswordResetMessageTemplate)
	if err != nil {
		log.Fatal(err)
	}
	return &passwordResetOutputHandler{
		emailClient: emailClient,
		template:    template,
	}
}
