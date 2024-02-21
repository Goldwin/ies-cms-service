package queries

import "github.com/Goldwin/ies-pik-cms/pkg/common/queries"

type ViewPersonByEmailQuery struct {
	Email string
}

type ViewPersonByEmail interface {
	Execute(query ViewPersonByEmailQuery) (ViewPersonResult, queries.QueryErrorDetail)
}
