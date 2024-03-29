package common

import "github.com/Goldwin/ies-pik-cms/pkg/common/out"

type NoopOutput[T any] struct{}

// OnError implements out.Output.
func (*NoopOutput[T]) OnError(err out.AppErrorDetail) {
	//noop
}

// OnSuccess implements out.Output.
func (*NoopOutput[T]) OnSuccess(result T) {
	//noop
}
