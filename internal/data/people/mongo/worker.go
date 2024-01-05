package mongo

import (
	"context"

	"github.com/Goldwin/ies-pik-cms/pkg/common/worker"
	"github.com/Goldwin/ies-pik-cms/pkg/people/commands"
	"github.com/Goldwin/ies-pik-cms/pkg/people/queries"
	"go.mongodb.org/mongo-driver/mongo"
)

type unitOfWorkImpl struct {
	mongoClient    *mongo.Client
	useTransaction bool
	dbName         string
}

type queryWorkerImpl struct {
	mongoClient    *mongo.Client
	dbName         string
	useTransaction bool
}

// Query implements worker.QueryWorker.
func (q *queryWorkerImpl) Query(ctx context.Context) queries.QueryContext {
	return NewQueryContext(ctx, q.mongoClient.Database(q.dbName))
}

// Execute implements worker.UnitOfWork.
func (u *unitOfWorkImpl) Execute(ctx context.Context, op worker.AtomicOperation[commands.CommandContext]) error {
	db := u.mongoClient.Database(u.dbName)
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

func NewUnitOfWork(mongoClient *mongo.Client, databaseName string, useTransaction bool) worker.UnitOfWork[commands.CommandContext] {
	return &unitOfWorkImpl{
		mongoClient:    mongoClient,
		useTransaction: useTransaction,
		dbName:         databaseName,
	}
}

func NewQueryWorker(mongoClient *mongo.Client, databaseName string, useTransaction bool) worker.QueryWorker[queries.QueryContext] {
	return &queryWorkerImpl{
		mongoClient:    mongoClient,
		useTransaction: useTransaction,
		dbName:         databaseName,
	}
}
