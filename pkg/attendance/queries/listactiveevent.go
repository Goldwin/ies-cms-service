package queries

import (
	"context"
	"time"

	"github.com/Goldwin/ies-pik-cms/pkg/attendance/dto"
)

type ListActiveEventFilter struct {
	FromDate time.Time
	ToDate   time.Time
}

type ListActiveEventResult struct {
	Data []dto.EventDTO
}

type ListActiveEvent interface {
	Execute(ctx context.Context, filter ListActiveEventFilter) (*ListActiveEventResult, error)
}
