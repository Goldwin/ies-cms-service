package infra

import (
	"fmt"
	"net/smtp"
)

type EmailConfig struct {
	Host     string `env:"EMAIL_HOST" yaml:"host" default:"localhost"`
	Port     int    `env:"EMAIL_PORT" yaml:"port" default:"25"`
	Username string `env:"EMAIL_USERNAME" yaml:"username" default:""`
	Password string `env:"EMAIL_PASSWORD" yaml:"password" default:""`
}

type emailClientImpl struct {
	config EmailConfig
}

type EmailClient interface {
	Send(to []string, subject string, mime string, message string) error
}

func NewEmailClient(config EmailConfig) EmailClient {
	return &emailClientImpl{
		config: config,
	}
}

func (s *emailClientImpl) Send(to []string, subject string, mime string, message string) error {

	hostAndPort := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)

	auth := smtp.PlainAuth("", s.config.Username, s.config.Password, s.config.Host)

	msg := fmt.Sprintf("Subject: %s\n%s\n%s", subject, mime, message)
	return smtp.SendMail(hostAndPort, auth, s.config.Username, to, []byte(msg))
}
