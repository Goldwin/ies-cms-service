package repositories

import (
	"github.com/Goldwin/ies-pik-cms/pkg/auth/entities"
	"github.com/Goldwin/ies-pik-cms/pkg/common/repositories"
)

type OtpRepository interface {
	repositories.Repository[string, entities.Otp]
}
