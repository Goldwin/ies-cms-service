package queries

import (
	"github.com/Goldwin/ies-pik-cms/pkg/people/dto"
)

type ViewPersonQuery struct {
	ID string
}

type ViewPersonResult struct {
	Data dto.Person `json:"data"`
}

type ViewPerson interface {
	Execute(query ViewPersonQuery) (ViewPersonResult, error)
}
