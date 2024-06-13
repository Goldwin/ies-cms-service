package mongo

import (
	"context"

	"github.com/Goldwin/ies-pik-cms/pkg/attendance/commands"
	"github.com/Goldwin/ies-pik-cms/pkg/common/worker"
)

type attendanceUnitOfWorkImpl struct {
}

// Execute implements worker.UnitOfWork.
func (a *attendanceUnitOfWorkImpl) Execute(ctx context.Context, op worker.AtomicOperation[commands.CommandContext]) error {
	panic("unimplemented")
}

func NewUnitOfWork() worker.UnitOfWork[commands.CommandContext] {
	return &attendanceUnitOfWorkImpl{}
}
