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

type attendanceRepositoryImpl struct {
	ctx        context.Context
	db         *mongo.Database
	collection *mongo.Collection
}

// Delete implements repositories.EventScheduleRepository.
func (e *attendanceRepositoryImpl) Delete(attendance *entities.Attendance) error {
	_, err := e.collection.DeleteOne(e.ctx, bson.M{"_id": attendance.ID})
	return err
}

// Get implements repositories.EventScheduleRepository.
func (e *attendanceRepositoryImpl) Get(id string) (*entities.Attendance, error) {
	var model AttendanceModel
	err := e.collection.FindOne(e.ctx, bson.M{"_id": id}).Decode(&model)

	if err == mongo.ErrNoDocuments {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return model.ToAttendance(), nil
}

// List implements repositories.EventScheduleRepository.
func (e *attendanceRepositoryImpl) List(idList []string) ([]*entities.Attendance, error) {
	var models []AttendanceModel
	cursor, err := e.collection.Find(e.ctx, bson.M{"_id": bson.M{"$in": idList}})

	if err != nil {
		return nil, err
	}

	defer cursor.Close(e.ctx)
	if err = cursor.All(e.ctx, &models); err != nil {
		return nil, err
	}

	return lo.Map(models, func(model AttendanceModel, _ int) *entities.Attendance {
		return model.ToAttendance()
	}), nil
}

// Save implements repositories.EventScheduleRepository.
func (e *attendanceRepositoryImpl) Save(attendance *entities.Attendance) (*entities.Attendance, error) {
	model := toAttendanceModel(attendance)

	_, err := e.collection.UpdateByID(e.ctx, attendance.ID, bson.M{"$set": model}, options.Update().SetUpsert(true))
	if err != nil {
		return nil, err
	}
	return attendance, nil
}

func NewAttendanceRepository(ctx context.Context, db *mongo.Database) repositories.AttendanceRepository {
	return &attendanceRepositoryImpl{
		ctx:        ctx,
		db:         db,
		collection: db.Collection(AttendanceCollection),
	}
}
