package queries

import (
	"github.com/Goldwin/ies-pik-cms/pkg/common/queries"
	"github.com/Goldwin/ies-pik-cms/pkg/people/dto"
)

type ViewPersonQuery struct {
	ID string
}

type ViewPersonResult queries.QueryResult[dto.Person]

type ViewPerson interface {
	Execute(query ViewPersonQuery) (ViewPersonResult, queries.QueryErrorDetail)
}
