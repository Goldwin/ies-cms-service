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

type PasswordModel struct {
	EmailAddress string `bson:"_id"`
	Salt         []byte `bson:"salt"`
	PasswordHash []byte `bson:"passwordHash"`
}

func (m *PasswordModel) toEntity() *entities.PasswordDetail {
	return &entities.PasswordDetail{
		EmailAddress: entities.EmailAddress(m.EmailAddress),
		Salt:         m.Salt,
		PasswordHash: m.PasswordHash,
	}
}

func fromPasswordDetailEntity(e *entities.PasswordDetail) PasswordModel {
	return PasswordModel{
		EmailAddress: string(e.EmailAddress),
		Salt:         e.Salt,
		PasswordHash: e.PasswordHash,
	}
}
