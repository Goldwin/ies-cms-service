package controllers

import (
	"github.com/Goldwin/ies-pik-cms/pkg/common/out"
)

type outputDecorator[T any] struct {
	output      out.Output[T]
	errFunction func(out.AppErrorDetail)
	successFunc func(T)
}

// OnError implements out.Output.
func (o *outputDecorator[T]) OnError(err out.AppErrorDetail) {
	if o.output != nil {
		o.output.OnError(err)
	}
	if o.errFunction != nil {
		o.errFunction(err)
	}
}

// OnSuccess implements out.Output.
func (o *outputDecorator[T]) OnSuccess(result T) {
	if o.output != nil {
		o.output.OnSuccess(result)
	}
	if o.successFunc != nil {
		o.successFunc(result)
	}
}
