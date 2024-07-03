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
func (l *listEventAttendanceImpl) Execute(query ListEventAttendanceQuery) (ListEventAttendanceResult, queries.QueryErrorDetail) {
	filter := bson.M{"eventId": bson.M{"$eq": query.EventID}}

	if query.EventActivityID != "" {
		filter["activityId"] = bson.M{"$eq": query.EventActivityID}
	}

	if len(query.AttendanceTypes) > 0 {
		filter["type"] = bson.M{"$in": query.AttendanceTypes}
	}

	cursor, err := l.db.Collection(AttendanceCollection).Find(
		l.ctx, filter,
		options.Find().SetSort(bson.D{{Key: "checkinTime", Value: 1}}),
		options.Find().SetLimit(int64(query.Limit)),
	)

	if err != nil {
		return ListEventAttendanceResult{}, queries.QueryErrorDetail{
			Code:    500,
			Message: "Failed to connect to database",
		}
	}
	defer cursor.Close(l.ctx)
	var result = make([]AttendanceModel, 0)
	err = cursor.Decode(&result)

	if err != nil {
		return ListEventAttendanceResult{}, queries.QueryErrorDetail{
			Code:    500,
			Message: "Failed to Decode Event Information",
		}
	}

	activityCache, err := l.fetchActivities(query)

	if err != nil {
		return ListEventAttendanceResult{}, queries.QueryErrorDetail{
			Code:    500,
			Message: "Failed to Fetch Event Activities Information",
		}
	}

	return ListEventAttendanceResult{
		Data: lo.Map(result, func(e AttendanceModel, _ int) dto.EventAttendanceDTO {
			return dto.EventAttendanceDTO{
				ID:                e.ID,
				EventID:           e.EventID,
				Activity:          activityCache[e.EventActivityID],
				PersonID:          e.PersonID,
				FirstName:         e.FirstName,
				MiddleName:        e.MiddleName,
				LastName:          e.LastName,
				ProfilePictureURL: e.ProfilePictureUrl,
				SecurityCode:      e.SecurityCode,
				SecurityNumber:    e.SecurityNumber,
				CheckinTime:       e.CheckinTime,
				AttendanceType:    e.Type,
			}
		}),
	}, queries.NoQueryError
}

func (l *listEventAttendanceImpl) fetchActivities(query ListEventAttendanceQuery) (activityCache map[string]dto.EventActivityDTO, err error) {
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
