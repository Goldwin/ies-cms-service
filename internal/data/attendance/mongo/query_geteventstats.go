package mongo

import (
	"context"
	"time"

	"github.com/Goldwin/ies-pik-cms/pkg/attendance/dto"
	. "github.com/Goldwin/ies-pik-cms/pkg/attendance/queries"
	"github.com/Goldwin/ies-pik-cms/pkg/common/queries"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type getEventScheduleStatsImpl struct {
	db  *mongo.Database
	ctx context.Context
}

// Execute implements queries.GetEventScheduleStats.
func (g *getEventScheduleStatsImpl) Execute(filter GetEventScheduleStatsFilter) (GetEventScheduleStatsResult, queries.QueryErrorDetail) {
	var models []EventAttendanceSummaryModel
	schedule, err := g.getSchedule(filter.ScheduleID)

	if err == mongo.ErrNoDocuments {
		return GetEventScheduleStatsResult{}, queries.ResourceNotFoundError("Schedule not found")
	}
	if err != nil {
		return GetEventScheduleStatsResult{}, queries.InternalServerError(err)
	}

	coll := g.db.Collection(AttendanceSummaryCollection)
	findOption := options.Find().SetSort(bson.D{bson.E{Key: "date", Value: 1}}).
		SetProjection(bson.D{
			bson.E{Key: "_id", Value: 1},
			bson.E{Key: "scheduleId", Value: 1},
			bson.E{Key: "date", Value: 1},
			bson.E{Key: "totalByType", Value: 1},
		}).SetLimit(20)

	monthsAgo := time.Now().AddDate(0, -1, 0)
	cursor, err := coll.Find(g.ctx, bson.M{"scheduleId": filter.ScheduleID, "date": bson.M{"$gte": monthsAgo, "$lte": time.Now()}}, findOption)

	if err != nil {
		return GetEventScheduleStatsResult{}, queries.InternalServerError(err)
	}
	defer cursor.Close(g.ctx)

	err = cursor.All(g.ctx, &models)

	if err != nil {
		return GetEventScheduleStatsResult{}, queries.InternalServerError(err)
	}

	result := dto.EventScheduleStatsDTO{
		ID:         filter.ScheduleID,
		EventStats: make([]dto.EventStatsDTO, 0),
	}

	for _, model := range models {
		result.EventStats = append(result.EventStats, dto.EventStatsDTO{
			ID:   model.ID,
			Date: model.Date.Add(time.Duration(schedule.TimezoneOffset) * time.Hour).Format("02 Jan 2006"),
			AttendanceCount: lo.MapToSlice(model.TotalByType, func(k string, v int) dto.EventAttendanceCountStats {
				return dto.EventAttendanceCountStats{
					AttendanceType: k,
					Count:          v,
				}
			}),
		})
	}

	return GetEventScheduleStatsResult{
		Data: result,
	}, queries.NoQueryError
}

func (g *getEventScheduleStatsImpl) getSchedule(scheduleID string) (EventScheduleModel, error) {
	var schedule EventScheduleModel
	coll := g.db.Collection(EventScheduleCollection)
	err := coll.FindOne(g.ctx, bson.M{"_id": scheduleID}).Decode(&schedule)
	if err != nil {
		return EventScheduleModel{}, err
	}
	return schedule, nil
}

func NewGetEventScheduleStats(ctx context.Context, db *mongo.Database) GetEventScheduleStats {
	return &getEventScheduleStatsImpl{
		db:  db,
		ctx: ctx,
	}
}
