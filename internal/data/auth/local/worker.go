package local

import (
	"context"

	"github.com/Goldwin/ies-pik-cms/pkg/auth/repositories"
	"github.com/Goldwin/ies-pik-cms/pkg/common/worker"
)

type authUnitOfWorkImpl struct {
}

func (u *authUnitOfWorkImpl) Execute(ctx context.Context, op worker.AtomicOperation[repositories.CommandContext]) error {
	return op(NewContext())
}

func NewUnitOfWork() worker.UnitOfWork[repositories.CommandContext] {
	return &authUnitOfWorkImpl{}
}
