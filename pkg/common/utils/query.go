package utils

import (
	"sync"

	"github.com/Goldwin/ies-pik-cms/pkg/common/out"
	"github.com/Goldwin/ies-pik-cms/pkg/common/queries"
)

type Filter interface {
	Validate() error
}

type Query[T Filter, R any] interface {
	Execute(T) (R, queries.QueryErrorDetail)
}

type singleQueryExecutor[T Filter, R any] struct {
	query  Query[T, R]
	output out.Output[R]
}

func SingleQueryExecution[T Filter, R any](
	query Query[T, R]) *singleQueryExecutor[T, R] {
	return &singleQueryExecutor[T, R]{
		query:  query,
		output: &noopOutput[R]{},
	}
}

func (s *singleQueryExecutor[T, R]) WithOutput(output out.Output[R]) *singleQueryExecutor[T, R] {
	s.output = output
	return s
}

func (s *singleQueryExecutor[T, R]) Execute(filter T) out.Waitable {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := filter.Validate(); err != nil {
			s.output.OnError(out.AppErrorDetail{
				Code:    400,
				Message: err.Error(),
			})			
			return
		}
		result, err := s.query.Execute(filter)
		if err.NoError() {
			s.output.OnSuccess(result)
		} else {
			s.output.OnError(out.ConvertQueryErrorDetail(err))
		}
	}()
	return wg
}
