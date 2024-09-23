package repositories

import (
	"github.com/Goldwin/ies-pik-cms/pkg/auth/entities"
	"github.com/Goldwin/ies-pik-cms/pkg/common/repositories"
)

type PasswordRepository interface {
	repositories.Repository[string, entities.PasswordDetail]
}
