package commands

import (
	"fmt"

	. "github.com/Goldwin/ies-pik-cms/pkg/common/commands"
	"github.com/Goldwin/ies-pik-cms/pkg/people/dto"
	"github.com/Goldwin/ies-pik-cms/pkg/people/entities"
	"github.com/Goldwin/ies-pik-cms/pkg/people/repositories"
)

/*
Add Person use case, where it will create a new person's account, and link inputted emails with user account
*/
type UpdatePersonCommand struct {
	Input dto.Person
}

const (
	UpdatePersonErrorCodeUserNotExist AppErrorCode = 10102
	UpdatePersonErrorCodeEmailsExist  AppErrorCode = 10102
	UpdatePersonErrorCodeInvalidInput AppErrorCode = 10103
)

func (cmd UpdatePersonCommand) Execute(ctx repositories.CommandContext) AppExecutionResult[dto.Person] {
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

	personResult, err := ctx.PersonRepository().Get(c.ID)

	if personResult == nil {
		return AppExecutionResult[dto.Person]{
			Status: ExecutionStatusFailed,
			Error: AppErrorDetail{
				Code:    UpdatePersonErrorCodeUserNotExist,
				Message: fmt.Sprintf("Can't Update Person Info, Error: Person with id %s does not exist", c.ID),
			},
		}
	}

	person := entities.Person{
		ID:                c.ID,
		FirstName:         c.FirstName,
		MiddleName:        c.MiddleName,
		LastName:          c.LastName,
		Addresses:         addresses,
		PhoneNumbers:      phones,
		ProfilePictureUrl: c.ProfilePictureUrl,
		EmailAddress:      entities.EmailAddress(c.EmailAddress),
		MaritalStatus:     c.MaritalStatus,
		Birthday:          c.Birthday.ToEntity(),
		Gender:            entities.Gender(c.Gender),
	}

	err = person.Validate()

	if err != nil {
		return AppExecutionResult[dto.Person]{
			Status: ExecutionStatusFailed,
			Error: AppErrorDetail{
				Code:    AddPersonErrorCodeInvalidInput,
				Message: fmt.Sprintf("Can't Update Person Info, Error: %s", err.Error()),
			},
		}
	}

	result, err := ctx.PersonRepository().UpdatePerson(person)

	if err != nil {
		return AppExecutionResult[dto.Person]{
			Status: ExecutionStatusFailed,
			Error: AppErrorDetail{
				Code:    AddHouseholdErrorCodeDBError,
				Message: fmt.Sprintf("Can't Add Update Person Info, Error: %s", err.Error()),
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
