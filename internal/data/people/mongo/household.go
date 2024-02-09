package mongo

import (
	"context"

	"github.com/Goldwin/ies-pik-cms/pkg/people/entities"
	"github.com/Goldwin/ies-pik-cms/pkg/people/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type householdRepositoryImpl struct {
	ctx                       context.Context
	db                        *mongo.Database
	householdCollection       *mongo.Collection
	personHouseholdCollection *mongo.Collection
}

// GetHousehold implements repositories.HouseholdRepository.
func (h *householdRepositoryImpl) GetHousehold(id string) (*entities.Household, error) {
	var household Household
	err := h.householdCollection.FindOne(h.ctx, bson.M{"_id": id}).Decode(&household)
	if err != nil {
		return nil, err
	}
	entities := toHouseholdEntities(household)
	return &entities, nil
}

// AddHousehold implements repositories.HouseholdRepository.
func (h *householdRepositoryImpl) AddHousehold(e entities.Household) (*entities.Household, error) {
	//share with head's id
	e.ID = e.HouseholdHead.ID
	_, err := h.householdCollection.InsertOne(h.ctx, toHouseholdModel(e))

	totalMembers := len(e.Members) + 1
	personIds := make([]string, totalMembers)
	for i, member := range e.Members {
		personIds[i] = member.ID
	}
	personIds[totalMembers-1] = e.HouseholdHead.ID

	h.personHouseholdCollection.UpdateMany(h.ctx,
		bson.M{"personID": bson.M{"$in": personIds}},
		bson.M{"$set": bson.M{"householdID": e.ID}},
		options.Update().SetUpsert(true),
	)

	if err != nil {
		return nil, err
	}
	return &e, nil
}

// UpdateHousehold implements repositories.HouseholdRepository.
func (h *householdRepositoryImpl) UpdateHousehold(e entities.Household) (*entities.Household, error) {
	var err error
	newHousehold := toHouseholdModel(e)
	newHousehold.ID = newHousehold.HouseholdHead.PersonID

	//replace old household's member id to no household
	_, err = h.personHouseholdCollection.UpdateMany(h.ctx, bson.M{"householdID": e.ID}, bson.M{"$set": bson.M{"householdID": ""}})

	if err != nil {
		return nil, err
	}

	if e.ID == newHousehold.ID {
		_, err = h.householdCollection.UpdateOne(h.ctx, bson.M{"_id": e.ID}, bson.M{"$set": toHouseholdModel(e)})
	} else {
		_, err = h.householdCollection.DeleteOne(h.ctx, bson.M{"_id": e.ID})
		if err != nil {
			return nil, err
		}
		_, err = h.householdCollection.InsertOne(h.ctx, newHousehold)

		if err != nil {
			return nil, err
		}
	}

	totalMembers := len(e.Members) + 1

	personIds := make([]string, totalMembers)
	for i, member := range e.Members {
		personIds[i] = member.ID
	}

	personIds[totalMembers-1] = e.HouseholdHead.ID

	//replace member's household id with new ids
	_, err = h.personHouseholdCollection.UpdateMany(h.ctx,
		bson.M{"personID": bson.M{"$in": personIds}},
		bson.M{"$set": bson.M{"householdID": newHousehold.ID}},
		options.Update().SetUpsert(true),
	)

	if err != nil {
		return nil, err
	}

	return &e, nil
}

func NewHouseholdRepository(ctx context.Context, db *mongo.Database) repositories.HouseholdRepository {
	return &householdRepositoryImpl{
		ctx:                       ctx,
		db:                        db,
		householdCollection:       db.Collection("household"),
		personHouseholdCollection: db.Collection("personHousehold"),
	}
}
