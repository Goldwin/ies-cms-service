package local

import (
	"context"

	"github.com/Goldwin/ies-pik-cms/pkg/auth/commands"
	"github.com/Goldwin/ies-pik-cms/pkg/common/worker"
)

type authUnitOfWorkImpl struct {
}

func (u *authUnitOfWorkImpl) Execute(ctx context.Context, op worker.AtomicOperation[commands.CommandContext]) error {
	return op(NewContext())
}

func NewUnitOfWork() worker.UnitOfWork[commands.CommandContext] {
	return &authUnitOfWorkImpl{}
}
