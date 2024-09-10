package mongo

import (
	"context"

	"github.com/Goldwin/ies-pik-cms/pkg/attendance/dto"
	. "github.com/Goldwin/ies-pik-cms/pkg/attendance/queries"
	"github.com/Goldwin/ies-pik-cms/pkg/common/queries"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type listEventAttendanceImpl struct {
	db  *mongo.Database
	ctx context.Context
}

// Execute implements queries.ListEventAttendance.
func (l *listEventAttendanceImpl) Execute(query ListEventAttendanceFilter) (ListEventAttendanceResult, queries.QueryErrorDetail) {
	filter := bson.M{"event._id": bson.M{"$eq": query.EventID}}

	if query.EventActivityID != "" {
		filter["activityId"] = bson.M{"$eq": query.EventActivityID}
	}

	if len(query.AttendanceTypes) > 0 {
		filter["type"] = bson.M{"$in": query.AttendanceTypes}
	}

	if query.LastID != "" {
		filter["_id"] = bson.M{"$gt": query.LastID}
	}

	cursor, err := l.db.Collection(AttendanceCollection).Find(
		l.ctx, filter,
		options.Find().SetLimit(int64(query.Limit)).SetSort(bson.D{{Key: "_id", Value: 1}}),
	)

	if err != nil {
		return ListEventAttendanceResult{}, queries.InternalServerError(err)
	}

	defer cursor.Close(l.ctx)
	var result = make([]AttendanceModel, 0)
	err = cursor.All(l.ctx, &result)

	if err != nil {
		return ListEventAttendanceResult{}, queries.InternalServerError(err)
	}

	return ListEventAttendanceResult{
		Data: lo.Map(result, func(e AttendanceModel, _ int) dto.EventAttendanceDTO {
			return dto.EventAttendanceDTO{
				ID: e.ID,
				Event: dto.EventDTO{
					ID:         e.Event.ID,
					ScheduleID: e.Event.ScheduleID,
					Name:       e.Event.Name,
					Date:       e.Event.Date,
				},
				Activity: dto.EventActivityDTO{
					ID:   e.EventActivity.ID,
					Name: e.EventActivity.Name,
					Time: e.EventActivity.Time,
				},
				Attendee: dto.AttendeeDTO{
					PersonID:          e.Attendee.PersonID,
					FirstName:         e.Attendee.FirstName,
					MiddleName:        e.Attendee.MiddleName,
					LastName:          e.Attendee.LastName,
					ProfilePictureURL: e.Attendee.ProfilePictureUrl,
				},

				CheckedInBy: dto.AttendeeDTO{
					PersonID:          e.CheckedInBy.PersonID,
					FirstName:         e.CheckedInBy.FirstName,
					MiddleName:        e.CheckedInBy.MiddleName,
					LastName:          e.CheckedInBy.LastName,
					ProfilePictureURL: e.CheckedInBy.ProfilePictureUrl,
				},

				SecurityCode:   e.SecurityCode,
				SecurityNumber: e.SecurityNumber,
				CheckinTime:    e.CheckinTime,
				AttendanceType: e.Type,
			}
		}),
	}, queries.NoQueryError
}

func (l *listEventAttendanceImpl) fetchActivities(query ListEventAttendanceFilter) (activityCache map[string]dto.EventActivityDTO, err error) {
	var eventModel EventModel
	err = l.db.Collection(query.EventID).FindOne(l.ctx, bson.M{"id": query.EventActivityID}).Decode(&eventModel)
	if err != nil {
		return
	}
	activityCache = lo.SliceToMap(eventModel.EventActivities, func(e EventActivityModel) (string, dto.EventActivityDTO) {
		return e.ID, dto.EventActivityDTO{
			ID:   e.ID,
			Name: e.Name,
			Time: e.Time,
		}
	})
	return activityCache, nil
}

func NewListEventAttendance(ctx context.Context, db *mongo.Database) ListEventAttendance {
	return &listEventAttendanceImpl{
		db:  db,
		ctx: ctx,
	}
}
