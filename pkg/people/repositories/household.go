package repositories

import (
	"github.com/Goldwin/ies-pik-cms/pkg/common/repositories"
	"github.com/Goldwin/ies-pik-cms/pkg/people/entities"
)

type HouseholdRepository interface {		
	UpdateHousehold(*entities.Household) (*entities.Household, error)

	repositories.Repository[string, entities.Household]
}
