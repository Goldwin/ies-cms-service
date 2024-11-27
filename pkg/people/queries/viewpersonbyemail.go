package queries

import "github.com/Goldwin/ies-pik-cms/pkg/common/queries"

type ViewPersonByEmailFilter struct {
	Email string
}

func (q ViewPersonByEmailFilter) Validate() error {
	return nil
}

type ViewPersonByEmail interface {
	Execute(query ViewPersonByEmailFilter) (ViewPersonResult, queries.QueryErrorDetail)
}
