package repositories

import "github.com/Goldwin/ies-pik-cms/pkg/auth/entities"

type RoleRepository interface {
	Save(role *entities.Role) error
}
