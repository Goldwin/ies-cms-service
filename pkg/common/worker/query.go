package worker

import "context"

type QueryWorker[T any] interface {
	Query(context.Context) T
}
