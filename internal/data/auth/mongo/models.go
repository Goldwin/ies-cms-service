package mongo

import (
	"time"

	"github.com/Goldwin/ies-pik-cms/pkg/auth/entities"
	"github.com/samber/lo"
)

const (
	OTPCollection               = "otp"
	AccountCollection           = "accounts"
	PasswordCollection          = "passwords"
	ResetPasswordCodeCollection = "reset_password_codes"
)

type PasswordResetCodeModel struct {
	Email     string    `bson:"_id"`
	Code      string    `bson:"code"`
	ExpiresAt time.Time `bson:"expiryAt"`
}

func (m *PasswordResetCodeModel) toEntity() *entities.PasswordResetCode {
	return &entities.PasswordResetCode{
		Email:     m.Email,
		Code:      m.Code,
		ExpiresAt: m.ExpiresAt,
	}
}

func toPasswordResetCodeModel(e *entities.PasswordResetCode) PasswordResetCodeModel {
	return PasswordResetCodeModel{
		Email:     e.Email,
		Code:      e.Code,
		ExpiresAt: e.ExpiresAt,
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

type OTPModel struct {
	EmailAddress string `bson:"_id"`
	PasswordHash []byte `bson:"passwordHash"`
	Salt         []byte `bson:"salt"`
	ExpiresAt    time.Time
}

func fromOtpEntity(e *entities.Otp) OTPModel {
	return OTPModel{
		EmailAddress: string(e.EmailAddress),
		PasswordHash: e.PasswordHash,
		Salt:         e.Salt,
		ExpiresAt:    e.ExpiresAt,
	}
}

func (m *OTPModel) toEntity() *entities.Otp {
	return &entities.Otp{
		EmailAddress: entities.EmailAddress(m.EmailAddress),
		PasswordHash: m.PasswordHash,
		Salt:         m.Salt,
		ExpiresAt:    m.ExpiresAt,
	}
}

type AccountModel struct {
	Email string      `bson:"_id"`
	Roles []RoleModel `bson:"roles"`
}

func (m *AccountModel) toEntity() *entities.Account {
	return &entities.Account{
		Email: entities.EmailAddress(m.Email),
		Roles: lo.Map(m.Roles, func(role RoleModel, _ int) *entities.Role {
			return role.toEntity()
		}),
	}
}

func fromAccountEntity(e *entities.Account) AccountModel {
	return AccountModel{
		Email: string(e.Email),
		Roles: lo.Map(e.Roles, func(role *entities.Role, _ int) RoleModel {
			return fromRoleEntity(role)
		}),
	}
}

type RoleModel struct {
	ID     string   `bson:"_id"`
	Name   string   `bson:"name"`
	Scopes []string `bson:"scopes"`
}

func fromRoleEntity(e *entities.Role) RoleModel {
	return RoleModel{
		ID:   e.ID,
		Name: e.Name,
		Scopes: lo.Map(e.Scopes, func(scope entities.Scope, _ int) string {
			return string(scope)
		}),
	}
}

func (m *RoleModel) toEntity() *entities.Role {
	return &entities.Role{
		ID:   m.ID,
		Name: m.Name,
		Scopes: lo.Map(m.Scopes, func(scope string, _ int) entities.Scope {
			return entities.Scope(scope)
		}),
	}
}
