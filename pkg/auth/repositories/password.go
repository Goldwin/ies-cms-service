package repositories

import (
	"time"

	"github.com/Goldwin/ies-pik-cms/pkg/auth/entities"
	"github.com/Goldwin/ies-pik-cms/pkg/common/repositories"
)

type PasswordRepository interface {
	GetResetCode(e entities.EmailAddress) (string, error)
	SaveResetCode(e entities.EmailAddress, token string, ttl time.Duration) error
	DeleteResetToken(e entities.EmailAddress) error
	repositories.Repository[string, entities.PasswordDetail]
}
