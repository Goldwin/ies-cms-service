package worker

import (
	"context"
)

type AtomicOperation[CTX any] func(CTX) error

type UnitOfWork[CTX any] interface {
	Execute(ctx context.Context, op AtomicOperation[CTX]) error
}
