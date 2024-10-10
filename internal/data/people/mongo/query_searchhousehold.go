package mongo

import (
	"context"
	"fmt"
	"log"
	"time"

	. "github.com/Goldwin/ies-pik-cms/pkg/common/queries"
	"github.com/Goldwin/ies-pik-cms/pkg/people/dto"
	"github.com/Goldwin/ies-pik-cms/pkg/people/queries"
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
func (s *searchHouseholdImpl) Execute(filter queries.SearchHouseholdFilter) (queries.SearchHouseholdResult, QueryErrorDetail) {
	var filteredPerson []PersonModel
	var householdIdList []string
	var err QueryErrorDetail

	filteredPerson, err = s.queryPersonModel(filter)

	if !err.NoError() {
		return queries.SearchHouseholdResult{}, err
	}

	householdIdList, err = s.listHouseholdID(filteredPerson)
	if !err.NoError() {
		return queries.SearchHouseholdResult{}, err
	}

	result, err := s.listHousehold(householdIdList)

	return queries.SearchHouseholdResult{
		Data: result,
	}, err
}

func (s *searchHouseholdImpl) listHousehold(householdIdList []string) ([]dto.Household, QueryErrorDetail) {
	var filteredPersonHousehold []HouseholdModel
	householdCollection := s.db.Collection(householdCollectionName)

	householdCursor, err := householdCollection.Find(s.ctx, bson.M{"_id": bson.M{"$in": lo.Map(householdIdList, func(e string, _ int) string {
		return e
	})}})

	if err != nil {
		log.Default().Printf("Failed to query list of household: %s", err.Error())
		return nil, QueryErrorDetail{
			Code:    500,
			Message: "Failed to query list of household",
		}
	}
	defer householdCursor.Close(s.ctx)

	err = householdCursor.All(s.ctx, &filteredPersonHousehold)

	if err != nil {
		log.Default().Printf("Failed to query list of household: %s", err.Error())
		return nil, QueryErrorDetail{
			Code:    500,
			Message: "Failed to query list of household",
		}
	}

	householdMembersByHouseholdID, queryError := s.listPersonsGroupByHouseholdID(
		lo.Map(filteredPersonHousehold, func(e HouseholdModel, _ int) string {
			return e.ID
		})...,
	)

	if !queryError.NoError() {
		return nil, queryError
	}

	return lo.Map(filteredPersonHousehold, func(household HouseholdModel, _ int) dto.Household {
		members := lo.Filter(householdMembersByHouseholdID[household.ID], func(e PersonModel, _ int) bool {
			return household.HouseholdHeadID != e.ID
		})

		head := lo.FindOrElse(householdMembersByHouseholdID[household.ID], PersonModel{}, func(e PersonModel) bool {
			return e.ID == household.HouseholdHeadID
		})

		return dto.Household{
			ID:         household.ID,
			Name:       household.Name,
			PictureUrl: household.PictureUrl,
			Members: lo.Map(members, func(e PersonModel, _ int) dto.HouseholdPerson {
				return dto.HouseholdPerson{
					ID:                e.ID,
					FirstName:         e.FirstName,
					MiddleName:        e.MiddleName,
					LastName:          e.LastName,
					PhoneNumber:       e.PhoneNumber,
					EmailAddress:      e.EmailAddress,
					ProfilePictureUrl: e.ProfilePictureUrl,
					Birthday:          parseBirthdayString(e.Birthday),
				}
			}),
			HouseholdHead: dto.HouseholdPerson{
				ID:                head.ID,
				FirstName:         head.FirstName,
				MiddleName:        head.MiddleName,
				LastName:          head.LastName,
				PhoneNumber:       head.PhoneNumber,
				EmailAddress:      head.EmailAddress,
				ProfilePictureUrl: head.ProfilePictureUrl,
				Birthday:          parseBirthdayString(head.Birthday),
			},
		}
	}), NoQueryError
}

func parseBirthdayString(birthday *string) time.Time {
	if birthday == nil {
		return time.Time{}
	}
	var year, month, day int
	fmt.Sscanf(*birthday, "%d-%d-%d", &year, &month, &day)

	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}

func (s *searchHouseholdImpl) listPersonsGroupByHouseholdID(householdIDList ...string) (map[string][]PersonModel, QueryErrorDetail) {
	const errorMessage = "Failed to query list of person"
	cursor, err := s.db.Collection(personCollectionName).Find(s.ctx, bson.M{"householdId": bson.M{"$in": householdIDList}})

	if err != nil {
		log.Default().Printf("%s: %s", errorMessage, err.Error())
		return nil, QueryErrorDetail{
			Code:    500,
			Message: errorMessage,
		}
	}
	var personList []PersonModel
	err = cursor.All(s.ctx, &personList)

	if err != nil {
		log.Default().Printf("%s: %s", errorMessage, err.Error())
		return nil, QueryErrorDetail{
			Code:    500,
			Message: errorMessage,
		}
	}

	return lo.GroupBy(personList, func(model PersonModel) string {
		return model.HouseholdID
	}), NoQueryError
}

func (s *searchHouseholdImpl) listHouseholdID(person []PersonModel) ([]string, QueryErrorDetail) {
	return lo.Map(person, func(e PersonModel, _ int) string {
		return e.HouseholdID
	}), NoQueryError
}

func (s *searchHouseholdImpl) queryPersonModel(filter queries.SearchHouseholdFilter) ([]PersonModel, QueryErrorDetail) {
	var filteredPerson []PersonModel
	personCollection := s.db.Collection(personCollectionName)
	opts := options.Find().SetSort(bson.D{{Key: "_id", Value: 1}}).SetLimit(int64(filter.Limit)).SetProjection(bson.M{"_id": 1, "householdId": 1})
	regexOp := "$regex"
	opt := "$options"
	cursor, err := personCollection.Find(s.ctx,
		bson.M{
			"_id": bson.M{"$gt": filter.LastID},
			"$or": []interface{}{
				bson.M{"firstName": bson.M{regexOp: fmt.Sprintf("^%s", filter.NamePrefix), opt: "i"}},
				bson.M{"middleName": bson.M{regexOp: fmt.Sprintf("^%s", filter.NamePrefix), opt: "i"}},
				bson.M{"lastName": bson.M{regexOp: fmt.Sprintf("^%s", filter.NamePrefix), opt: "i"}},
			},
		},
		opts,
	)
	if err != nil {
		log.Default().Printf("Failed to query list of person: %s", err.Error())
		return nil, QueryErrorDetail{
			Code:    500,
			Message: "Failed to query list of person",
		}
	}
	defer cursor.Close(s.ctx)

	if err = cursor.All(s.ctx, &filteredPerson); err != nil {
		return nil, QueryErrorDetail{
			Code:    500,
			Message: "Failed to Decode Person Information",
		}
	}
	return filteredPerson, NoQueryError
}

func NewSearchHousehold(ctx context.Context, db *mongo.Database) queries.SearchHousehold {
	return &searchHouseholdImpl{
		ctx: ctx,
		db:  db,
	}
}
