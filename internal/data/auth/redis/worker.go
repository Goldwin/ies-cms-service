package redis

import (
	"context"

	"github.com/Goldwin/ies-pik-cms/pkg/auth/commands"
	"github.com/Goldwin/ies-pik-cms/pkg/common/worker"
	"github.com/redis/go-redis/v9"
)

type authUnitOfWorkImpl struct {
	client redis.UniversalClient
}

func (u *authUnitOfWorkImpl) Execute(ctx context.Context, op worker.AtomicOperation[commands.CommandContext]) error {
	pipe := u.client.TxPipeline()
	err := op(NewContext(ctx, u.client, pipe))
	if err != nil {
		pipe.Discard()
		return err
	}
	_, err = pipe.Exec(ctx)
	return err
}

func NewUnitOfWork(r redis.UniversalClient) worker.UnitOfWork[commands.CommandContext] {
	return &authUnitOfWorkImpl{
		client: r,
	}
}
