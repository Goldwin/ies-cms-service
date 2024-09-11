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

type getEventScheduleStatsImpl struct {
	db  *mongo.Database
	ctx context.Context
}

// Execute implements queries.GetEventScheduleStats.
func (g *getEventScheduleStatsImpl) Execute(filter GetEventScheduleStatsFilter) (GetEventScheduleStatsResult, queries.QueryErrorDetail) {
	var models []EventAttendanceSummaryModel
	coll := g.db.Collection(AttendanceSummaryCollection)
	findOption := options.Find().SetSort(bson.D{bson.E{Key: "date", Value: 1}}).
		SetProjection(bson.D{
			bson.E{Key: "_id", Value: 1},
			bson.E{Key: "scheduleId", Value: 1},
			bson.E{Key: "date", Value: 1},
			bson.E{Key: "totalByType", Value: 1},
		})
	cursor, err := coll.Find(g.ctx, bson.M{"scheduleId": filter.ScheduleID}, findOption)

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
			Date: model.Date,
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

func NewGetEventScheduleStats(ctx context.Context, db *mongo.Database) GetEventScheduleStats {
	return &getEventScheduleStatsImpl{
		db:  db,
		ctx: ctx,
	}
}
