package repositories

import (
	"time"

	"github.com/Goldwin/ies-pik-cms/pkg/auth/entities"
)

type PasswordRepository interface {
	Get(entities.EmailAddress) (*entities.PasswordDetail, error)
	Save(entities.PasswordDetail) error
	GetResetCode(e entities.EmailAddress) (string, error)
	SaveResetCode(e entities.EmailAddress, token string, ttl time.Duration) error
	DeleteResetToken(e entities.EmailAddress) error
}
