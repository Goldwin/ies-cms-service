package mongo

import (
	"context"
	"fmt"
	"log"

	"github.com/Goldwin/ies-pik-cms/pkg/attendance/dto"
	. "github.com/Goldwin/ies-pik-cms/pkg/attendance/queries"
	"github.com/Goldwin/ies-pik-cms/pkg/common/queries"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type searchHouseholdImpl struct {
	ctx context.Context
	db  *mongo.Database
}

// Execute implements queries.SearchHousehold.
func (s *searchHouseholdImpl) Execute(filter SearchHouseholdFilter) (SearchHouseholdResult, queries.QueryErrorDetail) {
	var filteredPerson []PersonModel
	var householdIdList []string
	var err queries.QueryErrorDetail

	filteredPerson, err = s.queryPersonModel(filter)

	if !err.NoError() {
		return SearchHouseholdResult{}, err
	}

	householdIdList, err = s.getHouseholdID(filteredPerson)
	if !err.NoError() {
		return SearchHouseholdResult{}, err
	}

	result, err := s.listHousehold(householdIdList)

	return SearchHouseholdResult{
		Data: result,
	}, err
}

func (s *searchHouseholdImpl) listHousehold(householdIdList []string) ([]dto.HouseholdDTO, queries.QueryErrorDetail) {
	var filteredPersonHousehold []HouseholdModel
	householdCollection := s.db.Collection(HouseholdCollection)

	householdCursor, err := householdCollection.Find(s.ctx, bson.M{"_id": bson.M{"$in": lo.Map(householdIdList, func(e string, _ int) string {
		return e
	})}})

	if err != nil {
		log.Default().Printf("Failed to query list of household: %s", err.Error())
		return nil, queries.QueryErrorDetail{
			Code:    500,
			Message: "Failed to query list of household",
		}
	}
	defer householdCursor.Close(s.ctx)

	householdCursor.All(s.ctx, &filteredPersonHousehold)

	return lo.Map(filteredPersonHousehold, func(e HouseholdModel, _ int) dto.HouseholdDTO {
		return dto.HouseholdDTO{
			ID:            e.ID,
			Name:          e.Name,
			PictureUrl:    e.PictureUrl,
			HouseholdHead: e.HouseholdHead.ToDTO(),
			Members: lo.Map(e.HouseholdMembers, func(e PersonModel, _ int) dto.PersonDTO {
				return e.ToDTO()
			}),
		}
	}), queries.NoQueryError
}

func (s *searchHouseholdImpl) getHouseholdID(person []PersonModel) ([]string, queries.QueryErrorDetail) {
	var filteredPersonHousehold []PersonHouseholdModel
	personHouseholdCollection := s.db.Collection(PersonHouseholdCollection)

	householdCursor, err := personHouseholdCollection.Find(s.ctx, bson.M{"_id": bson.M{"$in": lo.Map(person, func(e PersonModel, _ int) string {
		return e.ID
	})}})

	if err != nil {
		log.Default().Printf("Failed to query list of household ID: %s", err.Error())
		return nil, queries.QueryErrorDetail{
			Code:    500,
			Message: "Failed to query list of household ID",
		}
	}
	defer householdCursor.Close(s.ctx)
	if err = householdCursor.All(s.ctx, &filteredPersonHousehold); err != nil {
		return nil, queries.QueryErrorDetail{
			Code:    500,
			Message: "Failed to Decode Person Information",
		}
	}

	return lo.Map(filteredPersonHousehold, func(e PersonHouseholdModel, _ int) string {
		return e.HouseholdID
	}), queries.NoQueryError
}

func (s *searchHouseholdImpl) queryPersonModel(filter SearchHouseholdFilter) ([]PersonModel, queries.QueryErrorDetail) {
	var filteredPerson []PersonModel
	personCollection := s.db.Collection(PersonCollection)
	opts := options.Find().SetSort(bson.D{{Key: "_id", Value: 1}}).SetLimit(int64(filter.Limit)).SetProjection(bson.M{"_id": 1})
	regexOp := "$regex"
	cursor, err := personCollection.Find(s.ctx,
		bson.M{
			"_id": bson.M{"$gt": filter.LastID},
			"$or": []interface{}{
				bson.M{"firstName": bson.M{regexOp: fmt.Sprintf("^%s", filter.NamePrefix)}},
				bson.M{"middleName": bson.M{regexOp: fmt.Sprintf("^%s", filter.NamePrefix)}},
				bson.M{"lastName": bson.M{regexOp: fmt.Sprintf("^%s", filter.NamePrefix)}},
			},
		},
		opts,
	)
	if err != nil {
		log.Default().Printf("Failed to query list of person: %s", err.Error())
		return nil, queries.QueryErrorDetail{
			Code:    500,
			Message: "Failed to query list of person",
		}
	}
	defer cursor.Close(s.ctx)

	if err = cursor.All(s.ctx, &filteredPerson); err != nil {
		return nil, queries.QueryErrorDetail{
			Code:    500,
			Message: "Failed to Decode Person Information",
		}
	}
	return filteredPerson, queries.NoQueryError
}

func NewSearchHousehold(ctx context.Context, db *mongo.Database) SearchHousehold {
	return &searchHouseholdImpl{
		ctx: ctx,
		db:  db,
	}
}
