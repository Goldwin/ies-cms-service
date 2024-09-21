package queries

import (
	"github.com/Goldwin/ies-pik-cms/pkg/common/queries"
	"github.com/Goldwin/ies-pik-cms/pkg/people/dto"
)

type SearchPersonQuery struct {
	LastID     string `json:"lastId"`
	NamePrefix string `json:"namePrefix"`
	Limit      int    `json:"limit"`
}

type SearchPersonResult struct {
	Data []dto.Person `json:"data"`
}

type SearchPerson interface {
	Execute(query SearchPersonQuery) (SearchPersonResult, queries.QueryErrorDetail)
}
