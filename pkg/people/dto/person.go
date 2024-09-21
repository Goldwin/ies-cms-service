package dto

import (
	"fmt"

	"github.com/Goldwin/ies-pik-cms/pkg/people/entities"
)

type Address struct {
	Line1      string `json:"line1"`
	Line2      string `json:"line2"`
	City       string `json:"city"`
	Province   string `json:"province"`
	PostalCode string `json:"postalCode"`
}

type YearMonthDay string

type Person struct {
	ID                string        `json:"id"`
	FirstName         string        `json:"firstName"`
	MiddleName        string        `json:"middleName"`
	LastName          string        `json:"lastName"`
	ProfilePictureUrl string        `json:"profilePictureUrl"`
	Address           string        `json:"address"`
	PhoneNumber       string        `json:"phoneNumber"`
	EmailAddress      string        `json:"emailAddress"`
	MaritalStatus     string        `json:"maritalStatus"`
	Birthday          *YearMonthDay `json:"birthday"`
	Gender            string        `json:"gender"`
}

func (y *YearMonthDay) ToEntity() *entities.YearMonthDay {
	if y == nil || *y == "" {
		return nil
	}
	var ymd entities.YearMonthDay
	fmt.Sscanf(string(*y), "%d-%d-%d", &ymd.Year, &ymd.Month, &ymd.Day)
	return &ymd
}
