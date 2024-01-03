package commands

import (
	"fmt"

	. "github.com/Goldwin/ies-pik-cms/pkg/common/commands"
	"github.com/Goldwin/ies-pik-cms/pkg/people/dto"
	"github.com/Goldwin/ies-pik-cms/pkg/people/entities"
)

/*
Add Person use case, where it will create a new person's account, and link inputted emails with user account
*/
type AddPersonCommand struct {
	Input dto.Person
}

const (
	AddPersonErrorCodeEmailsExist  AppErrorCode = 10002
	AddPersonErrorCodeInvalidInput AppErrorCode = 10003
)

func (cmd AddPersonCommand) Execute(ctx CommandContext) AppExecutionResult[dto.Person] {
	var err error
	c := cmd.Input
	addresses := make([]entities.Address, len(c.Addresses))
	for i, address := range c.Addresses {
		addresses[i] = entities.Address(address)
	}

	phones := make([]entities.PhoneNumber, len(c.PhoneNumbers))
	for i, phone := range c.PhoneNumbers {
		phones[i] = entities.PhoneNumber(phone)
	}

	person := entities.Person{
		ID:                c.ID,
		FirstName:         c.FirstName,
		MiddleName:        c.MiddleName,
		LastName:          c.LastName,
		Addresses:         addresses,
		PhoneNumbers:      phones,
		EmailAddress:      entities.EmailAddress(c.EmailAddress),
		MaritalStatus:     c.MaritalStatus,
		ProfilePictureUrl: c.ProfilePictureUrl,
		Birthday:          c.Birthday.ToEntity(),
		Gender:            entities.Gender(c.Gender),
	}

	err = person.Validate()

	if err != nil {
		return AppExecutionResult[dto.Person]{
			Status: ExecutionStatusFailed,
			Error: AppErrorDetail{
				Code:    AddPersonErrorCodeInvalidInput,
				Message: fmt.Sprintf("Can't Add New Person Info, Error: %s", err.Error()),
			},
		}
	}

	//To be determined later, whether we need to validate email or not.
	// accounts, err := ctx.AccountRepository().FindAccountByEmails(emails)

	// if err != nil {
	// 	return AppExecutionResult[dto.Person]{
	// 		Status: ExecutionStatusFailed,
	// 		Error: AppErrorDetail{
	// 			Code:    AddHouseholdErrorCodeDBError,
	// 			Message: fmt.Sprintf("Can't Add New Person Info, Error: %s", err.Error()),
	// 		},
	// 	}
	// }

	// if len(accounts) > 0 {
	// 	return AppExecutionResult[dto.Person]{
	// 		Status: ExecutionStatusFailed,
	// 		Error: AppErrorDetail{
	// 			Code:    AddPersonErrorCodeEmailsExist,
	// 			Message: "Can't Add New Person Info, Email already exists in the database",
	// 		},
	// 	}
	// }

	result, err := ctx.PersonRepository().AddPerson(person)

	if err != nil {
		return AppExecutionResult[dto.Person]{
			Status: ExecutionStatusFailed,
			Error: AppErrorDetail{
				Code:    AddHouseholdErrorCodeDBError,
				Message: fmt.Sprintf("Can't Add New Person Info, Error: %s", err.Error()),
			},
		}
	}

	output := cmd.Input
	output.ID = result.ID

	return AppExecutionResult[dto.Person]{
		Status: ExecutionStatusSuccess,
		Error:  AppErrorDetailNone,
		Result: output,
	}
}
