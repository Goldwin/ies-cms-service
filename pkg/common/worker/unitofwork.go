package worker

import (
	"context"
)

type AtomicOperation[CTX any] func(CTX) error

/*
Unit of work is an abstraction to isolate all of state changes that happen during a business transaction
*/
type UnitOfWork[CTX any] interface {
	Execute(ctx context.Context, op AtomicOperation[CTX]) error
}
