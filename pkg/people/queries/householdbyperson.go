package queries

import (
	"github.com/Goldwin/ies-pik-cms/pkg/common/queries"
	"github.com/Goldwin/ies-pik-cms/pkg/people/dto"
)

type ViewHouseholdByPersonQuery struct {
	PersonID string
}

type ViewHouseholdByPersonResult queries.QueryResult[dto.Household]

type ViewHouseholdByPerson interface {
	Execute(query ViewHouseholdByPersonQuery) (ViewHouseholdByPersonResult, queries.QueryErrorDetail)
}