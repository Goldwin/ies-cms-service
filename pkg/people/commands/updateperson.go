package commands

import (
	"fmt"

	. "github.com/Goldwin/ies-pik-cms/pkg/common/commands"
	"github.com/Goldwin/ies-pik-cms/pkg/people/dto"
	"github.com/Goldwin/ies-pik-cms/pkg/people/entities"
)

/*
Update Person use case, where it will update a person's information
*/
type UpdatePersonCommand struct {
	Input dto.Person
}

const (
	UpdatePersonErrorCodeUserNotExist CommandErrorCode = 10012
	UpdatePersonErrorCodeEmailsExist  CommandErrorCode = 10013
	UpdatePersonErrorCodeInvalidInput CommandErrorCode = 10014
	UpdatePersonErrorCodeDBError      CommandErrorCode = 10015
)

func (cmd UpdatePersonCommand) Execute(ctx CommandContext) CommandExecutionResult[dto.Person] {
	var err error
	c := cmd.Input

	phones := make([]entities.PhoneNumber, len(c.PhoneNumber))
	for i, phone := range c.PhoneNumber {
		phones[i] = entities.PhoneNumber(phone)
	}

	personResult, err := ctx.PersonRepository().Get(c.ID)

	if personResult == nil {
		return CommandExecutionResult[dto.Person]{
			Status: ExecutionStatusFailed,
			Error: CommandErrorDetail{
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
		Address:           c.Address,
		PhoneNumber:       entities.PhoneNumber(c.PhoneNumber),
		ProfilePictureUrl: c.ProfilePictureUrl,
		EmailAddress:      entities.EmailAddress(c.EmailAddress),
		MaritalStatus:     c.MaritalStatus,
		Birthday:          c.Birthday.ToEntity(),
		Gender:            entities.Gender(c.Gender),
	}

	err = person.Validate()

	if err != nil {
		return CommandExecutionResult[dto.Person]{
			Status: ExecutionStatusFailed,
			Error: CommandErrorDetail{
				Code:    AddPersonErrorCodeInvalidInput,
				Message: fmt.Sprintf("Can't Update Person Info, Error: %s", err.Error()),
			},
		}
	}

	isEmailExist, err := checkEmailExistence(ctx, person.EmailAddress)
	if err != nil {
		return CommandExecutionResult[dto.Person]{
			Status: ExecutionStatusFailed,
			Error: CommandErrorDetail{
				Code:    UpdatePersonErrorCodeDBError,
				Message: fmt.Sprintf("Can't Update Person Info, Failed on checking email existence. Error: %s", err.Error()),
			},
		}
	}

	if isEmailExist {
		return CommandExecutionResult[dto.Person]{
			Status: ExecutionStatusFailed,
			Error: CommandErrorDetail{
				Code:    UpdatePersonErrorCodeEmailsExist,
				Message: fmt.Sprintf("Can't Update Person Info, email has been used by someone else"),
			},
		}
	}

	result, err := ctx.PersonRepository().UpdatePerson(person)

	if err != nil {
		return CommandExecutionResult[dto.Person]{
			Status: ExecutionStatusFailed,
			Error: CommandErrorDetail{
				Code:    AddHouseholdErrorCodeDBError,
				Message: fmt.Sprintf("Can't Add Update Person Info, Error: %s", err.Error()),
			},
		}
	}

	output := cmd.Input
	output.ID = result.ID

	return CommandExecutionResult[dto.Person]{
		Status: ExecutionStatusSuccess,
		Error:  CommandErrorDetailNone,
		Result: output,
	}
}
