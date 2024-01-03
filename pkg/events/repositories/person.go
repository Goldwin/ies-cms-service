package repositories

import "github.com/Goldwin/ies-pik-cms/pkg/events/entities"

type PersonRepository interface {
	Get(string) (*entities.Person, error)
	GetByIds([]string) ([]entities.Person, error)
}
