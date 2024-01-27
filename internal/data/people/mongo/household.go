package mongo

import (
	"context"

	"github.com/Goldwin/ies-pik-cms/pkg/people/entities"
	"github.com/Goldwin/ies-pik-cms/pkg/people/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Household struct {
	ID               string            `bson:"_id"`
	Name             string            `bson:"name"`
	HouseholdHead    HouseholdMember   `bson:"householdHead"`
	PictureUrl       string            `bson:"pictureUrl"`
	HouseholdMembers []HouseholdMember `bson:"householdMembers"`
}

type HouseholdMember struct {
	PersonID          string `bson:"personID"`
	FirstName         string `bson:"firstName"`
	LastName          string `bson:"lastName"`
	ProfilePictureUrl string `bson:"profilePictureUrl"`
	Email             string `bson:"email"`
	PhoneNumber       string `bson:"phoneNumber"`
}

type PersonHousehold struct {
	ID          string `bson:"_id"`
	HouseholdID string `bson:"householdID"`
}

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

func toHouseholdEntities(householdModel Household) entities.Household {
	return entities.Household{
		ID:            householdModel.ID,
		Name:          householdModel.Name,
		HouseholdHead: toHouseholdMemberEntities(householdModel.HouseholdHead),
		PictureUrl:    householdModel.PictureUrl,
		Members:       getMembersEntities(householdModel),
	}
}

func toHouseholdMemberEntities(e HouseholdMember) entities.Person {
	return entities.Person{
		ID:                e.PersonID,
		FirstName:         e.FirstName,
		LastName:          e.LastName,
		ProfilePictureUrl: e.ProfilePictureUrl,
		EmailAddress:      entities.EmailAddress(e.Email),
		PhoneNumber:       entities.PhoneNumber(e.PhoneNumber),
	}
}

func toHouseholdModel(e entities.Household) *Household {
	householdMembers := getMembersModel(e)
	return &Household{
		ID:               e.ID,
		Name:             e.Name,
		HouseholdHead:    toHouseholdMemberModel(e.HouseholdHead),
		PictureUrl:       e.PictureUrl,
		HouseholdMembers: householdMembers,
	}
}

func toHouseholdMemberModel(e entities.Person) HouseholdMember {
	return HouseholdMember{
		PersonID:          e.ID,
		FirstName:         e.FirstName,
		LastName:          e.LastName,
		ProfilePictureUrl: e.ProfilePictureUrl,
		Email:             string(e.EmailAddress),
		PhoneNumber:       string(e.PhoneNumber),
	}
}

func getMembersEntities(e Household) []entities.Person {
	householdMembers := make([]entities.Person, 0)
	for _, member := range e.HouseholdMembers {
		householdMembers = append(householdMembers, toHouseholdMemberEntities(member))
	}
	return householdMembers
}

func getMembersModel(e entities.Household) []HouseholdMember {
	householdMembers := make([]HouseholdMember, 0)
	for _, member := range e.Members {
		householdMembers = append(householdMembers, toHouseholdMemberModel(member))
	}
	return householdMembers
}
