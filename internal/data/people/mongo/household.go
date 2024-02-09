package mongo

import (
	"context"

	"github.com/Goldwin/ies-pik-cms/pkg/people/entities"
	"github.com/Goldwin/ies-pik-cms/pkg/people/repositories"
	"github.com/google/uuid"
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
	e.ID = uuid.NewString()
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
	oldHousehold, err := h.GetHousehold(e.ID)

	if err != nil {
		return nil, err
	}

	oldMemberIdSet := make(map[string]bool, len(oldHousehold.Members)+1)
	for _, member := range oldHousehold.Members {
		oldMemberIdSet[member.ID] = true
	}
	oldMemberIdSet[oldHousehold.HouseholdHead.ID] = true

	if err != nil {
		return nil, err
	}

	_, err = h.householdCollection.UpdateOne(h.ctx, bson.M{"_id": e.ID}, bson.M{"$set": toHouseholdModel(e)})

	totalMembers := len(e.Members) + 1

	personIds := make([]string, totalMembers)
	oldPersonIds := make([]string, 0)
	for i, member := range e.Members {
		personIds[i] = member.ID
		if oldMemberIdSet[member.ID] {
			oldMemberIdSet[member.ID] = false
		}
	}

	for id, isDiscarded := range oldMemberIdSet {
		if isDiscarded {
			oldPersonIds = append(oldPersonIds, id)
		}
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

	_, err = h.personHouseholdCollection.UpdateMany(h.ctx,
		bson.M{"personID": bson.M{"$in": personIds}},
		bson.M{"$set": bson.M{"householdID": ""}},
		options.Update().SetUpsert(true),
	)

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
