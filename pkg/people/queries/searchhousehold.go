package queries

import (
	"fmt"

	"github.com/Goldwin/ies-pik-cms/pkg/people/dto"
	"github.com/Goldwin/ies-pik-cms/pkg/common/queries"
)

type SearchHouseholdFilter struct {
	Limit      int    `json:"limit" form:"limit"`
	LastID     string `json:"lastId" form:"lastId"`
	NamePrefix string `json:"namePrefix" form:"namePrefix"`
}

func (query SearchHouseholdFilter) Validate() error {
	if query.Limit > 200 {
		return fmt.Errorf("limit must be less than or equal to 200")
	}

	if query.Limit <= 0 {
		return fmt.Errorf("limit must be greater than 0")
	}

	if len(query.NamePrefix) > 1000 {
		return fmt.Errorf("name prefix must be less than or equal to 100")
	}
	return nil
}

type SearchHouseholdResult struct {
	Data []dto.Household `json:"data" form:"data"`
}

type SearchHousehold interface {
	Execute(filter SearchHouseholdFilter) (SearchHouseholdResult, queries.QueryErrorDetail)
}
