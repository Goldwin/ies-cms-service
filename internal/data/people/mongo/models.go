package mongo

import (
	"fmt"

	"github.com/Goldwin/ies-pik-cms/pkg/people/dto"
	"github.com/Goldwin/ies-pik-cms/pkg/people/entities"
)

const (
	personHouseholdCollectionName = "person_households"
	personCollectionName      = "persons"
	householdCollectionName   = "households"
)

type Person struct {
	ID                string  `bson:"_id"`
	FirstName         string  `bson:"firstName"`
	MiddleName        string  `bson:"middleName"`
	LastName          string  `bson:"lastName"`
	ProfilePictureUrl string  `bson:"profilePictureUrl"`
	Address           string  `bson:"address"`
	PhoneNumber       string  `bson:"phoneNumber"`
	EmailAddress      string  `bson:"emailAddress"`
	MaritalStatus     string  `bson:"maritalStatus"`
	Birthday          *string `bson:"birthday"`
	Gender            string  `bson:"gender"`
}

func toPersonMongoModel(e entities.Person) Person {
	var birthday *string
	if e.Birthday != nil {
		str := fmt.Sprintf("%04d-%02d-%02d", e.Birthday.Year, e.Birthday.Month, e.Birthday.Day)
		birthday = &str
	} else {
		birthday = nil
	}

	return Person{
		ID:                e.ID,
		FirstName:         e.FirstName,
		MiddleName:        e.MiddleName,
		LastName:          e.LastName,
		ProfilePictureUrl: e.ProfilePictureUrl,
		Address:           e.Address,
		PhoneNumber:       string(e.PhoneNumber),
		EmailAddress:      string(e.EmailAddress),
		MaritalStatus:     e.MaritalStatus,
		Birthday:          birthday,
		Gender:            string(e.Gender),
	}
}

func toPersonEntities(p Person) entities.Person {
	var birthday *entities.YearMonthDay

	phones := make([]entities.PhoneNumber, len(p.PhoneNumber))

	for i, phone := range p.PhoneNumber {
		phones[i] = entities.PhoneNumber(phone)
	}

	if p.Birthday == nil {
		birthday = nil
	} else {
		birthday = &entities.YearMonthDay{}
		fmt.Sscanf(*p.Birthday, "%d-%d-%d", &birthday.Year, &birthday.Month, &birthday.Day)
	}

	return entities.Person{
		ID:                p.ID,
		FirstName:         p.FirstName,
		MiddleName:        p.MiddleName,
		LastName:          p.LastName,
		ProfilePictureUrl: p.ProfilePictureUrl,
		Address:           p.Address,
		PhoneNumber:       entities.PhoneNumber(p.PhoneNumber),
		EmailAddress:      entities.EmailAddress(p.EmailAddress),
		MaritalStatus:     p.MaritalStatus,
		Birthday:          birthday,
		Gender:            entities.Gender(p.Gender),
	}
}

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

func toHouseholdMemberDto(e HouseholdMember) dto.HouseholdPerson {
	return dto.HouseholdPerson{
		ID:                e.PersonID,
		FirstName:         e.FirstName,
		LastName:          e.LastName,
		ProfilePictureUrl: e.ProfilePictureUrl,
		EmailAddress:      e.Email,
		PhoneNumber:       e.PhoneNumber,
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
