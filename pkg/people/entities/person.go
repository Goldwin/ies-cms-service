package entities

import (
	"fmt"
	"net/mail"
	"regexp"
	"time"
)

type PhoneNumber string

type EmailAddress string

type Gender string

const (
	Male   Gender = "MALE"
	Female Gender = "FEMALE"
)

type YearMonthDay struct {
	Year  int
	Month int
	Day   int
}

type Person struct {
	ID                string
	FirstName         string
	MiddleName        string
	LastName          string
	ProfilePictureUrl string
	Address           string
	PhoneNumber       PhoneNumber
	EmailAddress      EmailAddress
	MaritalStatus     string
	Birthday          *YearMonthDay
	Gender            Gender
}

func (p *Person) Validate() error {
	if p.FirstName == "" || p.LastName == "" {
		return fmt.Errorf("First name and Last name must be filled")
	}

	if !p.EmailAddress.IsValid() {
		return fmt.Errorf("invalid email: %s", p.EmailAddress)
	}

	if !p.PhoneNumber.IsValid() {
		return fmt.Errorf("invalid phone number: %s", p.PhoneNumber)
	}

	if p.Birthday != nil && !p.Birthday.IsValid() {
		return fmt.Errorf("invalid birthday: %04d-%02d-%02d", p.Birthday.Year, p.Birthday.Month, p.Birthday.Day)
	}

	return nil
}

func (e EmailAddress) IsValid() bool {
	if e == "" {
		return true
	}
	_, err := mail.ParseAddress(string(e))
	if err != nil {
		return false
	}
	return true
}

func (p PhoneNumber) IsValid() bool {
	if p == "" {
		return true
	}
	// Regular expression pattern for a phone number
	pattern := `^\+?[1-9]\d{1,14}$|^0[1-9]\d{1,14}$`

	regex := regexp.MustCompile(pattern)
	return regex.MatchString(string(p))
}

func (y YearMonthDay) IsValid() bool {
	_, err := time.Parse("2006-01-02", fmt.Sprintf("%04d-%02d-%02d", y.Year, y.Month, y.Day))
	return err == nil
}
