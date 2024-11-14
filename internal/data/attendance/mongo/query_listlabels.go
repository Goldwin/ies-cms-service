package mongo

import (
	"context"
	"log"

	"github.com/Goldwin/ies-pik-cms/pkg/attendance/dto"
	. "github.com/Goldwin/ies-pik-cms/pkg/attendance/queries"
	"github.com/Goldwin/ies-pik-cms/pkg/common/queries"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type listLabelsImpl struct {
	db     *mongo.Database
	ctx    context.Context
	logger *log.Logger
}

// Execute implements queries.ListLabels.
func (l *listLabelsImpl) Execute(filter ListLabelsFilter) (ListLabelsResult, queries.QueryErrorDetail) {
	var result []LabelModel
	cursor, err := l.db.Collection(LabelsCollection).Find(l.ctx, bson.M{"_id": bson.M{"$gt": filter.LastID}}, options.Find().SetLimit(int64(filter.Limit)))

	if err != nil {
		return ListLabelsResult{}, queries.InternalServerError(err)
	}
	defer cursor.Close(l.ctx)
	if err = cursor.All(l.ctx, &result); err != nil {
		return ListLabelsResult{}, queries.InternalServerError(err)
	}
	return ListLabelsResult{
		Data: lo.Map(result, func(model LabelModel, _ int) dto.LabelDTO {
			return model.ToDTO()
		}),
	}, queries.NoQueryError
}

func NewListLabels(ctx context.Context, db *mongo.Database) ListLabels {
	return &listLabelsImpl{
		ctx:    ctx,
		db:     db,
		logger: log.Default(),
	}
}
