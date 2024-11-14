package queries

import (
	"fmt"

	"github.com/Goldwin/ies-pik-cms/pkg/attendance/dto"
	"github.com/Goldwin/ies-pik-cms/pkg/common/queries"
)

type ListLabelsFilter struct {
	LastID string `json:"lastId" form:"lastId"`
	Limit  int    `json:"limit" form:"limit"`
}

func (query ListLabelsFilter) Validate() error {
	if query.Limit > 200 {
		return fmt.Errorf("limit must be less than or equal to 200")
	}

	if query.Limit <= 0 {
		return fmt.Errorf("limit must be greater than 0")
	}

	return nil
}

type ListLabelsResult struct {
	Data []dto.LabelDTO `json:"data"`
}

type ListLabels interface {
	Execute(ListLabelsFilter) (ListLabelsResult, queries.QueryErrorDetail)
}
