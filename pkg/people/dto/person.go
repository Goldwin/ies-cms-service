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

type YearMonthDay string

type Person struct {
	ID                string        `json:"id"`
	FirstName         string        `json:"first_name"`
	MiddleName        string        `json:"middle_name"`
	LastName          string        `json:"last_name"`
	ProfilePictureUrl string        `json:"profile_picture_url"`
	Address           string        `json:"address"`
	PhoneNumber       string        `json:"phone_number"`
	EmailAddress      string        `json:"email_address"`
	MaritalStatus     string        `json:"marital_status"`
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
