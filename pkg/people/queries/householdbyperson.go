package queries

import (
	"github.com/Goldwin/ies-pik-cms/pkg/common/queries"
	"github.com/Goldwin/ies-pik-cms/pkg/people/dto"
)

type ViewHouseholdByPersonFilter struct {
	PersonID string
}

func (q ViewHouseholdByPersonFilter) Validate() error {
	return nil
}

type ViewHouseholdByPersonResult queries.QueryResult[dto.Household]

type ViewHouseholdByPerson interface {
	Execute(query ViewHouseholdByPersonFilter) (ViewHouseholdByPersonResult, queries.QueryErrorDetail)
}
