package queries

import (
	"github.com/Goldwin/ies-pik-cms/pkg/people/dto"
)

//TODO: implement

type SearchPersonQuery struct {
	LastID string `json:"last_id"`
	Limit  int    `json:"limit"`
}

type SearchPersonResult struct {
	Data []dto.Person `json:"data"`
}

type SearchPerson interface {
	Execute(query SearchPersonQuery) (SearchPersonResult, error)
}
