package mongo

import (
	"time"

	"github.com/Goldwin/ies-pik-cms/pkg/auth/entities"
)

const (
	OTPCollection               = "otp"
	AccountCollection           = "account"
	PasswordCollection          = "password"
	ResetPasswordCodeCollection = "reset_password_code"
)

type PasswordResetCodeModel struct {
	Email    string    `bson:"_id"`
	Code     string    `bson:"code"`
	ExpiryAt time.Time `bson:"expiryAt"`
}

func (m *PasswordResetCodeModel) toEntity() *entities.PasswordResetCode {
	return &entities.PasswordResetCode{
		Email:    m.Email,
		Code:     m.Code,
		ExpiryAt: m.ExpiryAt,
	}
}

func toPasswordResetCodeModel(e *entities.PasswordResetCode) PasswordResetCodeModel {
	return PasswordResetCodeModel{
		Email:    e.Email,
		Code:     e.Code,
		ExpiryAt: e.ExpiryAt,
	}
}
