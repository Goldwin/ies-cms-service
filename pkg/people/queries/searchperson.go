package queries

import (
	"fmt"

	"github.com/Goldwin/ies-pik-cms/pkg/common/queries"
	"github.com/Goldwin/ies-pik-cms/pkg/people/dto"
)

type SearchPersonFilter struct {
	LastID     string `json:"lastId"`
	NamePrefix string `json:"namePrefix"`
	Limit      int    `json:"limit"`
}

func (query SearchPersonFilter) Validate() error {
	if query.Limit > 200 {
		return fmt.Errorf("limit must be less than or equal to 200")
	}

	if query.Limit <= 0 {
		return fmt.Errorf("limit must be greater than 0")
	}

	if len(query.NamePrefix) > 100 {
		return fmt.Errorf("name prefix must be less than or equal to 100")
	}

	return nil
}

type SearchPersonResult struct {
	Data []dto.Person `json:"data"`
}

type SearchPerson interface {
	Execute(query SearchPersonFilter) (SearchPersonResult, queries.QueryErrorDetail)
}
