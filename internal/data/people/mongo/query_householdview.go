package mongo

import (
	"context"
	"fmt"

	. "github.com/Goldwin/ies-pik-cms/pkg/common/queries"
	"github.com/Goldwin/ies-pik-cms/pkg/people/dto"
	"github.com/Goldwin/ies-pik-cms/pkg/people/queries"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type viewHouseholdByPersonImpl struct {
	ctx                 context.Context
	db                  *mongo.Database
	householdCollection *mongo.Collection
	personCollection    *mongo.Collection
}

const (
	connectionFailureMessage = "Failed to connect to database"
)

// Execute implements queries.ViewHouseholdByPerson.
func (v *viewHouseholdByPersonImpl) Execute(query queries.ViewHouseholdByPersonQuery) (queries.ViewHouseholdByPersonResult, QueryErrorDetail) {
	personModel, queryError := v.getPerson(query.PersonID)
	if queryError != NoQueryError {
		return queries.ViewHouseholdByPersonResult{}, queryError
	}

	household := HouseholdModel{}
	err := v.householdCollection.FindOne(v.ctx, bson.M{"_id": personModel.HouseholdID}).Decode(&household)
	if err != nil && err == mongo.ErrNoDocuments {
		return queries.ViewHouseholdByPersonResult{}, QueryErrorDetail{
			Code:    404,
			Message: fmt.Sprintf("Household Not Found"),
		}
	}

	householdMembers, queryError := v.listHouseholdMembers(household.ID)
	if queryError != NoQueryError {
		return queries.ViewHouseholdByPersonResult{}, queryError
	}

	householdHeadModel := lo.FindOrElse(householdMembers, PersonModel{}, func(e PersonModel) bool {
		return e.ID == household.HouseholdHeadID
	})

	householdMembers = lo.Filter(householdMembers, func(e PersonModel, _ int) bool {
		return e.ID != householdHeadModel.ID
	})

	return queries.ViewHouseholdByPersonResult{
		Data: &dto.Household{
			ID: household.ID,
			HouseholdHead: dto.HouseholdPerson{
				ID:                householdHeadModel.ID,
				FirstName:         householdHeadModel.FirstName,
				MiddleName:        householdHeadModel.MiddleName,
				LastName:          householdHeadModel.LastName,
				PhoneNumber:       householdHeadModel.PhoneNumber,
				EmailAddress:      householdHeadModel.EmailAddress,
				ProfilePictureUrl: householdHeadModel.ProfilePictureUrl,
			},
			Members: lo.Map(householdMembers, func(e PersonModel, _ int) dto.HouseholdPerson {
				return dto.HouseholdPerson{
					ID:                e.ID,
					FirstName:         e.FirstName,
					MiddleName:        e.MiddleName,
					LastName:          e.LastName,
					PhoneNumber:       e.PhoneNumber,
					EmailAddress:      e.EmailAddress,
					ProfilePictureUrl: e.ProfilePictureUrl,
				}
			}),
			Name: household.Name,
		},
	}, NoQueryError
}

func (v *viewHouseholdByPersonImpl) getPerson(personID string) (PersonModel, QueryErrorDetail) {
	var person PersonModel
	err := v.personCollection.FindOne(v.ctx, bson.M{"_id": personID}).Decode(&person)
	if err != nil && err == mongo.ErrNoDocuments {
		return PersonModel{}, QueryErrorDetail{
			Code:    404,
			Message: fmt.Sprintf("Person %s does not exists.", personID),
		}
	}
	if err != nil {
		return PersonModel{}, QueryErrorDetail{
			Code:    500,
			Message: connectionFailureMessage,
		}
	}
	return person, NoQueryError
}

func (v *viewHouseholdByPersonImpl) listHouseholdMembers(householdID string) ([]PersonModel, QueryErrorDetail) {
	var personList []PersonModel
	cursor, err := v.personCollection.Find(v.ctx, bson.M{"householdId": householdID})

	if err != nil {
		return nil, QueryErrorDetail{
			Code:    500,
			Message: connectionFailureMessage,
		}
	}

	if err = cursor.All(v.ctx, &personList); err != nil {
		return nil, QueryErrorDetail{
			Code:    500,
			Message: connectionFailureMessage,
		}
	}
	return personList, NoQueryError
}

func ViewHouseholdByPerson(ctx context.Context, db *mongo.Database) queries.ViewHouseholdByPerson {
	return &viewHouseholdByPersonImpl{
		ctx:                 ctx,
		db:                  db,
		householdCollection: db.Collection(householdCollectionName),
		personCollection:    db.Collection(personCollectionName),
	}
}
