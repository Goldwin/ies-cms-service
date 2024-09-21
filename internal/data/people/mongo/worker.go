package mongo

import (
	"context"

	"github.com/Goldwin/ies-pik-cms/pkg/common/worker"
	"github.com/Goldwin/ies-pik-cms/pkg/people/commands"
	"github.com/Goldwin/ies-pik-cms/pkg/people/queries"
	"go.mongodb.org/mongo-driver/mongo"
)

type unitOfWorkImpl struct {
	mongoDatabase  *mongo.Database
	useTransaction bool
	dbName         string
}

type queryWorkerImpl struct {
	mongoDatabase  *mongo.Database
	useTransaction bool
}

// Query implements worker.QueryWorker.
func (q *queryWorkerImpl) Query(ctx context.Context) queries.QueryContext {
	return NewQueryContext(ctx, q.mongoDatabase)
}

// Execute implements worker.UnitOfWork.
func (u *unitOfWorkImpl) Execute(ctx context.Context, op worker.AtomicOperation[commands.CommandContext]) error {
	db := u.mongoDatabase
	session, err := db.Client().StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	//transaction can only be done if mongo use replica set. otherwise it will fail.
	if u.useTransaction {
		_, err = session.WithTransaction(ctx, func(sessionContext mongo.SessionContext) (interface{}, error) {
			c := NewCommandContext(sessionContext, db)
			return nil, op(c)
		})
		return err
	}

	c := NewCommandContext(ctx, db)
	return op(c)
}

func NewUnitOfWork(mongoDatabase *mongo.Database, useTransaction bool) worker.UnitOfWork[commands.CommandContext] {
	return &unitOfWorkImpl{
		mongoDatabase:  mongoDatabase,
		useTransaction: useTransaction,
	}
}

func NewQueryWorker(mongoDatabase *mongo.Database, useTransaction bool) worker.QueryWorker[queries.QueryContext] {
	return &queryWorkerImpl{
		mongoDatabase:  mongoDatabase,
		useTransaction: useTransaction,
	}
}
