package utils

import (
	"context"
	"sync"

	"github.com/Goldwin/ies-pik-cms/pkg/common/commands"
	"github.com/Goldwin/ies-pik-cms/pkg/common/out"
	"github.com/Goldwin/ies-pik-cms/pkg/common/worker"
)

type singleCommandExecutor[CTX any, R any] struct {
	unitOfWork worker.UnitOfWork[CTX]
	command    commands.Command[CTX, R]

	output out.Output[R]
}

func (noopOutput[R]) OnError(err out.AppErrorDetail) {
	// noop
}

func (noopOutput[R]) OnSuccess(result R) {
	// noop
}

type noopOutput[R any] struct{}

func SingleCommandExecution[CTX any, R any](
	unitOfWork worker.UnitOfWork[CTX],
	command commands.Command[CTX, R]) *singleCommandExecutor[CTX, R] {
	return &singleCommandExecutor[CTX, R]{
		unitOfWork: unitOfWork,
		command:    command,
		output:     &noopOutput[R]{},
	}
}

func (s *singleCommandExecutor[CTX, R]) WithOutput(output out.Output[R]) *singleCommandExecutor[CTX, R] {
	s.output = output
	return s
}

func (s *singleCommandExecutor[CTX, R]) Execute(ctx context.Context) out.Waitable {
	wg := &sync.WaitGroup{}
	go s.unitOfWork.Execute(ctx, func(commandContext CTX) error {
		var res commands.CommandExecutionResult[R]
		wg.Add(1)
		defer wg.Done()
		res = s.command.Execute(commandContext)
		if res.Status == commands.ExecutionStatusSuccess {
			s.output.OnSuccess(res.Result)
			return nil
		}
		s.output.OnError(out.ConvertCommandErrorDetail(res.Error))
		return res.Error
	})
	return wg
}
