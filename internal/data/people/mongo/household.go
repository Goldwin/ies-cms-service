package mongo

import (
	"context"

	"github.com/Goldwin/ies-pik-cms/pkg/people/entities"
	"github.com/Goldwin/ies-pik-cms/pkg/people/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Household struct {
	ID               string            `bson:"_id"`
	Name             string            `bson:"name"`
	HouseholdHead    HouseholdMember   `bson:"householdHead"`
	PictureUrl       string            `bson:"pictureUrl"`
	HouseholdMembers []HouseholdMember `bson:"householdMembers"`
}

type HouseholdMember struct {
	MemberID          string `bson:"memberID"`
	FirstName         string `bson:"firstName"`
	LastName          string `bson:"lastName"`
	ProfilePictureUrl string `bson:"profilePictureUrl"`
	Email             string `bson:"email"`
	PhoneNumber       string `bson:"phoneNumber"`
}

type householdRepositoryImpl struct {
	ctx                 context.Context
	db                  *mongo.Database
	householdCollection *mongo.Collection
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
	if err != nil {
		return nil, err
	}
	return &e, nil
}

// UpdateHousehold implements repositories.HouseholdRepository.
func (h *householdRepositoryImpl) UpdateHousehold(e entities.Household) (*entities.Household, error) {
	var err error
	newHousehold := toHouseholdModel(e)
	newHousehold.ID = newHousehold.HouseholdHead.MemberID
	if e.ID == newHousehold.ID {
		_, err = h.householdCollection.UpdateOne(h.ctx, bson.M{"_id": e.ID}, bson.M{"$set": toHouseholdModel(e)})

	} else {
		_, err = h.householdCollection.DeleteOne(h.ctx, bson.M{"_id": e.ID})
		if err != nil {
			return nil, err
		}
		_, err = h.householdCollection.InsertOne(h.ctx, newHousehold)
	}
	if err != nil {
		return nil, err
	}

	return &e, nil
}

func NewHouseholdRepository(ctx context.Context, db *mongo.Database) repositories.HouseholdRepository {
	return &householdRepositoryImpl{
		ctx:                 ctx,
		db:                  db,
		householdCollection: db.Collection("household"),
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
		ID:                e.MemberID,
		FirstName:         e.FirstName,
		LastName:          e.LastName,
		ProfilePictureUrl: e.ProfilePictureUrl,
		EmailAddress:      entities.EmailAddress(e.Email),
		PhoneNumbers:      []entities.PhoneNumber{entities.PhoneNumber(e.PhoneNumber)},
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
		MemberID:          e.ID,
		FirstName:         e.FirstName,
		LastName:          e.LastName,
		ProfilePictureUrl: e.ProfilePictureUrl,
		Email:             string(e.EmailAddress),
		PhoneNumber:       string(e.PhoneNumbers[0]),
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
