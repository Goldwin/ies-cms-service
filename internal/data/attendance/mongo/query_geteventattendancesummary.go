package mongo

import (
	"context"
	"log"
	"time"

	. "github.com/Goldwin/ies-pik-cms/pkg/attendance/queries"
	"github.com/Goldwin/ies-pik-cms/pkg/common/queries"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type eventAttendanceAggregateKey struct {
	ActivityId string `bson:"activityId"`
	Name       string `bson:"name"`
	Type       string `bson:"type"`
}

type eventAttendanceAggregateModel struct {
	Key            eventAttendanceAggregateKey `bson:"_id"`
	Count          int                         `bson:"count"`
	FirstTimeCount int                         `bson:"firstTimeCount"`
}

type getEventAttendanceSummaryImpl struct {
	db  *mongo.Database
	ctx context.Context
}

// Execute implements queries.GetEventAttendanceSummary.
func (g *getEventAttendanceSummaryImpl) Execute(filter GetEventAttendanceSummaryFilter) (GetEventAttendanceSummaryResult, queries.QueryErrorDetail) {
	var summaryModel EventAttendanceSummaryModel
	var eventModel EventModel
	err := g.db.Collection(EventCollection).FindOne(g.ctx, bson.M{"_id": filter.EventID}).Decode(&eventModel)

	if err == mongo.ErrNoDocuments {
		return GetEventAttendanceSummaryResult{}, queries.ResourceNotFoundError("Event not found")
	}

	if err != nil {
		return GetEventAttendanceSummaryResult{}, queries.InternalServerError(err)
	}

	err = g.db.Collection(AttendanceSummaryCollection).FindOne(g.ctx, bson.M{"_id": filter.EventID}).Decode(&summaryModel)
	if err != nil && err != mongo.ErrNoDocuments {
		return GetEventAttendanceSummaryResult{}, queries.InternalServerError(err)
	}

	if err == mongo.ErrNoDocuments {
		summaryModel.ID = filter.EventID
		summaryModel.Date = eventModel.StartDate
		summaryModel.LastUpdated = time.Now().Add(-5 * time.Minute)
		summaryModel.TotalByType = map[string]int{}
	}

	go g.maybeUpdateSummary(summaryModel, eventModel)

	return GetEventAttendanceSummaryResult{
		Data: summaryModel.ToDTO(),
	}, queries.NoQueryError
}

func (g *getEventAttendanceSummaryImpl) maybeUpdateSummary(summary EventAttendanceSummaryModel, event EventModel) {
	if summary.NextUpdate.After(time.Now()) || summary.LastUpdated.Sub(event.StartDate) > 24*time.Hour {
		return
	}

	coll := g.db.Collection(AttendanceSummaryCollection)
	result := coll.FindOneAndUpdate(g.ctx,
		bson.M{"_id": summary.ID},
		bson.M{"$set": bson.M{"nextUpdate": time.Now().Add(5 * time.Minute)}},
		options.FindOneAndUpdate().SetUpsert(true),
	)

	err := result.Decode(&summary)

	if err != nil && err != mongo.ErrNoDocuments {
		log.Printf("Failed to update attendance summary. error: %s", err.Error())
		return
	}

	if err == mongo.ErrNoDocuments {
		summary.NextUpdate = time.Now().Add(-5 * time.Minute)
	}

	if summary.NextUpdate.After(time.Now()) {
		return
	}

	summary.ID = event.ID
	summary.Date = event.StartDate
	summary.ScheduleID = event.ScheduleID

	aggregates, err := g.aggregate(summary.ID)
	if err != nil {
		log.Printf("Failed to aggregate attendance. error: %s", err.Error())
		return
	}

	activitySummaries := make(map[string]ActivityAttendanceSummaryModel)
	summary.TotalByType = make(map[string]int)
	summary.Total = 0
	summary.TotalCheckedIn = 0
	summary.TotalCheckedOut = 0
	summary.TotalFirstTimer = 0

	for _, aggregate := range aggregates {
		activitySummary, ok := activitySummaries[aggregate.Key.ActivityId]
		if !ok {
			activitySummary = ActivityAttendanceSummaryModel{
				ID:          aggregate.Key.ActivityId,
				Name:        aggregate.Key.Name,
				Total:       0,
				TotalByType: make(map[string]int),
			}
		}

		activitySummary.Total += aggregate.Count
		activitySummary.TotalByType[aggregate.Key.Type] += aggregate.Count
		activitySummaries[aggregate.Key.ActivityId] = activitySummary

		summary.Total += aggregate.Count
		summary.TotalByType[aggregate.Key.Type] += aggregate.Count
		summary.TotalCheckedIn += aggregate.Count
		summary.TotalFirstTimer += aggregate.FirstTimeCount
	}

	summary.AcitivitiesSummary = make([]ActivityAttendanceSummaryModel, 0)

	for _, activitySummary := range activitySummaries {
		summary.AcitivitiesSummary = append(summary.AcitivitiesSummary, activitySummary)
	}

	summary.LastUpdated = time.Now()
	summary.NextUpdate = time.Now().Add(5 * time.Minute)
	_, err = coll.ReplaceOne(g.ctx, bson.M{"_id": summary.ID}, summary)

	if err != nil {
		log.Default().Printf("Failed to update summary, error: %s", err.Error())
	}
}

func (g *getEventAttendanceSummaryImpl) aggregate(eventID string) ([]eventAttendanceAggregateModel, error) {
	var models []eventAttendanceAggregateModel
	coll := g.db.Collection(AttendanceCollection)
	cursor, err := coll.Aggregate(g.ctx, mongo.Pipeline{
		bson.D{
			bson.E{
				Key: "$match",
				Value: bson.M{
					"event._id": bson.M{
						"$eq": eventID,
					},
				},
			},
		},
		bson.D{
			bson.E{
				Key: "$group",
				Value: bson.M{
					"_id": bson.M{
						"activityId": "$eventActivity._id",
						"name":       "$eventActivity.name",
						"type":       "$type",
					},
					"count": bson.M{
						"$sum": 1,
					},
					"firstTimeCount": bson.M{
						"$sum": bson.M{
							"$cond": bson.A{
								bson.M{
									"$eq": bson.A{
										"$firstTime", true,
									},
								},
								1, 0,
							},
						},
					},
				},
			},
		},
	})
	defer cursor.Close(g.ctx)

	if err != nil {
		return nil, err
	}
	cursor.All(g.ctx, &models)
	return models, nil
}

func NewGetEventAttendanceSummary(ctx context.Context, db *mongo.Database) GetEventAttendanceSummary {
	return &getEventAttendanceSummaryImpl{
		db:  db,
		ctx: ctx,
	}
}
