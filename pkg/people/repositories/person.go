package repositories

import "github.com/Goldwin/ies-pik-cms/pkg/people/entities"

type PersonRepository interface {
	AddPerson(entities.Person) (*entities.Person, error)
	UpdatePerson(entities.Person) (*entities.Person, error)
	Get(id string) (*entities.Person, error)
	ListByID(id []string) ([]entities.Person, error)
}
