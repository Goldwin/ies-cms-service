package entities

import "net/mail"

type EmailAddress string

type Person struct {
	ID         string
	FirstName  string
	MiddleName string
	LastName   string
}

type Account struct {
	Email  EmailAddress
	Roles  []Role
	Person Person
}

func (e EmailAddress) IsValid() bool {
	_, err := mail.ParseAddress(string(e))
	if err != nil {
		return false
	}
	return true
}
