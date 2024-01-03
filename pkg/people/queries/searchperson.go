package queries

import (
	"github.com/Goldwin/ies-pik-cms/pkg/people/dto"
)

//TODO: implement

type SearchPersonQuery struct {
	LastID string
	Limit  int
}

type SearchPersonResult struct {
	Data []dto.Person `json:"data"`
}

type SearchPerson interface {
	Execute(query SearchPersonQuery) (SearchPersonResult, error)
}
