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

type viewPersonImpl struct {
	ctx context.Context
	db  *mongo.Database
}

// Execute implements queries.ViewPerson.
func (v *viewPersonImpl) Execute(query queries.ViewPersonQuery) (queries.ViewPersonResult, QueryErrorDetail) {
	person := Person{}
	err := v.db.Collection("person").FindOne(v.ctx, bson.M{"_id": query.ID}).Decode(&person)
	if err != nil && err == mongo.ErrNoDocuments {
		return queries.ViewPersonResult{
				Data: nil,
			}, QueryErrorDetail{
				Code:    ResourceNotFound,
				Message: fmt.Sprintf("Person with id %s not found", query.ID),
			}
	}

	if err != nil {
		return queries.ViewPersonResult{}, QueryErrorDetail{
			Code:    500,
			Message: "Failed to connect to database",
		}
	}

	data := toPersonDTO(person)
	return queries.ViewPersonResult{
		Data: &data,
	}, QueryErrorDetail{}
}

func ViewPerson(ctx context.Context, db *mongo.Database) queries.ViewPerson {
	return &viewPersonImpl{
		ctx: ctx,
		db:  db,
	}
}

func toPersonDTO(person Person) dto.Person {

	result := dto.Person{
		ID:                person.ID,
		FirstName:         person.FirstName,
		MiddleName:        person.MiddleName,
		LastName:          person.LastName,
		ProfilePictureUrl: person.ProfilePictureUrl,
		Address:           person.Address,
		PhoneNumber:       person.PhoneNumber,
		EmailAddress:      person.EmailAddress,
		MaritalStatus:     person.MaritalStatus,
		Gender:            person.Gender,
	}

	if person.Birthday != nil {
		birthday := dto.YearMonthDay(*person.Birthday)
		result.Birthday = &birthday
	}
	return result
}
