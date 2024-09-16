package mongo

import (
	"context"

	"github.com/Goldwin/ies-pik-cms/pkg/attendance/entities"
	"github.com/Goldwin/ies-pik-cms/pkg/attendance/repositories"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type personAttendanceSummaryRepositoryImpl struct {
	ctx        context.Context
	db         *mongo.Database
	collection *mongo.Collection
}

// Delete implements repositories.PersonAttendanceSummaryRepository.
func (p *personAttendanceSummaryRepositoryImpl) Delete(e *entities.PersonAttendanceSummary) error {
	_, err := p.collection.DeleteOne(p.ctx, bson.M{"_id": e.PersonID})
	return err
}

// Get implements repositories.PersonAttendanceSummaryRepository.
func (p *personAttendanceSummaryRepositoryImpl) Get(id string) (*entities.PersonAttendanceSummary, error) {
	var model PersonAttendanceSummaryModel
	err := p.collection.FindOne(p.ctx, bson.M{"_id": id}).Decode(&model)

	if err == mongo.ErrNoDocuments {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return model.ToEntity(), nil
}

// List implements repositories.PersonAttendanceSummaryRepository.
func (p *personAttendanceSummaryRepositoryImpl) List(idList []string) ([]*entities.PersonAttendanceSummary, error) {
	var models []PersonAttendanceSummaryModel
	cursor, err := p.collection.Find(p.ctx, bson.M{"_id": bson.M{"$in": idList}})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(p.ctx)
	err = cursor.All(p.ctx, &models)
	if err != nil {
		return nil, err
	}
	return lo.Map(models, func(m PersonAttendanceSummaryModel, _ int) *entities.PersonAttendanceSummary {
		return m.ToEntity()
	}), nil
}

// Save implements repositories.PersonAttendanceSummaryRepository.
func (p *personAttendanceSummaryRepositoryImpl) Save(e *entities.PersonAttendanceSummary) (*entities.PersonAttendanceSummary, error) {
	model := toPersonAttendanceSummaryModel(e)
	_, err := p.collection.UpdateByID(p.ctx, model.ID, bson.M{"$set": model}, options.Update().SetUpsert(true))

	if err != nil {
		return nil, err
	}
	return e, nil
}

func NewPersonAttendanceSummaryRepository(ctx context.Context, db *mongo.Database) repositories.PersonAttendanceSummaryRepository {
	return &personAttendanceSummaryRepositoryImpl{
		ctx:        ctx,
		db:         db,
		collection: db.Collection(PersonAttendanceSummaryCollection),
	}
}
