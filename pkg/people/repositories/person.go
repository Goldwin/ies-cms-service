package repositories

import (
	"github.com/Goldwin/ies-pik-cms/pkg/common/repositories"
	"github.com/Goldwin/ies-pik-cms/pkg/people/entities"
)

type PersonRepository interface {			
	GetByEmail(email entities.EmailAddress) (*entities.Person, error)	

	repositories.Repository[string, entities.Person]
}
