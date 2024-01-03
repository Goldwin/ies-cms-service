package mongo

import (
	"context"

	"github.com/Goldwin/ies-pik-cms/pkg/common/worker"
	"github.com/Goldwin/ies-pik-cms/pkg/people/commands"
	"go.mongodb.org/mongo-driver/mongo"
)

type unitOfWorkImpl struct {
	mongoClient    *mongo.Client
	useTransaction bool
}

// Execute implements worker.UnitOfWork.
func (u *unitOfWorkImpl) Execute(ctx context.Context, op worker.AtomicOperation[commands.CommandContext]) error {
	db := u.mongoClient.Database("people")
	session, err := db.Client().StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	//transaction can only be done if mongo use replica set. otherwise it will fail.
	if u.useTransaction {
		_, err = session.WithTransaction(ctx, func(sessionContext mongo.SessionContext) (interface{}, error) {
			c := NewContext(sessionContext, db)
			return nil, op(c)
		})
		return err
	}

	c := NewContext(ctx, db)
	return op(c)
}

func NewUnitOfWork(mongoClient *mongo.Client, useTransaction bool) worker.UnitOfWork[commands.CommandContext] {
	return &unitOfWorkImpl{
		mongoClient:    mongoClient,
		useTransaction: useTransaction,
	}
}
