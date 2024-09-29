package repositories

import (
	"github.com/Goldwin/ies-pik-cms/pkg/common/repositories"
	"github.com/Goldwin/ies-pik-cms/pkg/people/entities"
)

type HouseholdRepository interface {		
	repositories.Repository[string, entities.Household]
}
