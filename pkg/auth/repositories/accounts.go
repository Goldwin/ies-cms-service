//go:generate mockery --output "mocks" --all --with-expecter=true
package repositories

import (
	"github.com/Goldwin/ies-pik-cms/pkg/auth/entities"
	"github.com/Goldwin/ies-pik-cms/pkg/common/repositories"
)

type AccountRepository interface {
	repositories.Repository[string, entities.Account]
}
