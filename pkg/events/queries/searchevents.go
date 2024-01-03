package queries

import (
	"github.com/Goldwin/ies-pik-cms/pkg/events/dto"
)

type SearchEventQuery struct {
	LastID string
	Limit  int
}

type SearchEventResult struct {
	Events []dto.ChurchEvent
}

type SearchEvent interface {
	Execute(query SearchEventQuery) (SearchEventResult, error)
}
