package mongo

import (
	"context"
	"fmt"

	. "github.com/Goldwin/ies-pik-cms/pkg/common/queries"
	"github.com/Goldwin/ies-pik-cms/pkg/people/dto"
	"github.com/Goldwin/ies-pik-cms/pkg/people/queries"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type viewHouseholdByPersonImpl struct {
	ctx                       context.Context
	db                        *mongo.Database
	householdCollection       *mongo.Collection
	personHouseholdCollection *mongo.Collection
}

// Execute implements queries.ViewHouseholdByPerson.
func (v *viewHouseholdByPersonImpl) Execute(query queries.ViewHouseholdByPersonQuery) (queries.ViewHouseholdByPersonResult, QueryErrorDetail) {
	personHousehold := PersonHouseholdModel{}
	err := v.personHouseholdCollection.FindOne(v.ctx, bson.M{"_id": query.PersonID}).Decode(&personHousehold)
	if err != nil && err == mongo.ErrNoDocuments {
		return queries.ViewHouseholdByPersonResult{}, QueryErrorDetail{
			Code:    404,
			Message: fmt.Sprintf("Person %s does not join any household yet.", query.PersonID),
		}
	}
	if err != nil {
		return queries.ViewHouseholdByPersonResult{}, QueryErrorDetail{
			Code:    500,
			Message: "Failed to connect to database",
		}
	}
	household := HouseholdModel{}
	err = v.householdCollection.FindOne(v.ctx, bson.M{"_id": personHousehold.HouseholdID}).Decode(&household)
	if err != nil && err == mongo.ErrNoDocuments {
		return queries.ViewHouseholdByPersonResult{}, QueryErrorDetail{
			Code:    404,
			Message: fmt.Sprintf("Household Not Found"),
		}
	}
	houseHoldMembers := make([]dto.HouseholdPerson, len(household.HouseholdMembers))

	for i, member := range household.HouseholdMembers {
		houseHoldMembers[i] = toHouseholdMemberDto(member)
	}
	return queries.ViewHouseholdByPersonResult{
		Data: &dto.Household{
			ID:            household.ID,
			HouseholdHead: toHouseholdMemberDto(household.HouseholdHead),
			Members:       houseHoldMembers,
			Name:          household.Name,
		},
	}, NoQueryError
}

func ViewHouseholdByPerson(ctx context.Context, db *mongo.Database) queries.ViewHouseholdByPerson {
	return &viewHouseholdByPersonImpl{
		ctx:                       ctx,
		db:                        db,
		householdCollection:       db.Collection(householdCollectionName),
		personHouseholdCollection: db.Collection(personHouseholdCollectionName),
	}
}
