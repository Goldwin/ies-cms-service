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
	PostalCode string `json:"postal_code"`
}

type PhoneNumber string

type EmailAddress string

type YearMonthDay string

type Person struct {
	ID                string        `json:"id"`
	FirstName         string        `json:"first_name"`
	MiddleName        string        `json:"middle_name"`
	LastName          string        `json:"last_name"`
	ProfilePictureUrl string        `json:"profile_picture_url"`
	Addresses         []Address     `json:"addresses"`
	PhoneNumbers      []PhoneNumber `json:"phone_numbers"`
	EmailAddress      EmailAddress  `json:"email_address"`
	MaritalStatus     string        `json:"marital_status"`
	Birthday          YearMonthDay  `json:"birthday"`
}

func (y YearMonthDay) ToEntity() entities.YearMonthDay {
	var ymd entities.YearMonthDay
	fmt.Sscanf(string(y), "%d-%d-%d", &ymd.Year, &ymd.Month, &ymd.Day)
	return ymd
}
