package repositories

import "github.com/Goldwin/ies-pik-cms/pkg/people/entities"

type HouseholdRepository interface {
	GetHousehold(string) (*entities.Household, error)
	AddHousehold(entities.Household) (*entities.Household, error)
	UpdateHousehold(entities.Household) (*entities.Household, error)
	DeleteHousehold(entities.Household) error
}
