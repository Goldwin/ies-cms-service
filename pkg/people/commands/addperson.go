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
	AddPersonErrorCodeEmailsExist  CommandErrorCode = 10002
	AddPersonErrorCodeInvalidInput CommandErrorCode = 10003
	AddPersonErrorCodeDBError      CommandErrorCode = 10004
)

func (cmd AddPersonCommand) Execute(ctx CommandContext) CommandExecutionResult[dto.Person] {
	var err error
	c := cmd.Input

	person := entities.Person{
		ID:                c.ID,
		FirstName:         c.FirstName,
		MiddleName:        c.MiddleName,
		LastName:          c.LastName,
		Address:           c.Address,
		PhoneNumber:       entities.PhoneNumber(c.PhoneNumber),
		EmailAddress:      entities.EmailAddress(c.EmailAddress),
		MaritalStatus:     c.MaritalStatus,
		ProfilePictureUrl: c.ProfilePictureUrl,
		Birthday:          c.Birthday.ToEntity(),
		Gender:            entities.Gender(c.Gender),
	}

	err = person.Validate()

	if err != nil {
		return CommandExecutionResult[dto.Person]{
			Status: ExecutionStatusFailed,
			Error: CommandErrorDetail{
				Code:    AddPersonErrorCodeInvalidInput,
				Message: fmt.Sprintf("Invalid Person Information, Error: %s", err.Error()),
			},
		}
	}

	//To be determined later, whether we need to validate email or not.

	isEmailExists, err := checkEmailExistence(ctx, person.EmailAddress)
	if err != nil {
		return CommandExecutionResult[dto.Person]{
			Status: ExecutionStatusFailed,
			Error: CommandErrorDetail{
				Code:    AddPersonErrorCodeDBError,
				Message: fmt.Sprintf("Can't Add New Person Info, Error: %s", err.Error()),
			},
		}
	}
	if isEmailExists {
		return CommandExecutionResult[dto.Person]{
			Status: ExecutionStatusFailed,
			Error: CommandErrorDetail{
				Code:    AddPersonErrorCodeEmailsExist,
				Message: "Can't Add New Person Info, Email already used by another person",
			},
		}
	}
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
		return CommandExecutionResult[dto.Person]{
			Status: ExecutionStatusFailed,
			Error: CommandErrorDetail{
				Code:    AddHouseholdErrorCodeDBError,
				Message: fmt.Sprintf("Can't Add New Person Info, Error: %s", err.Error()),
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

func checkEmailExistence(ctx CommandContext, emailAddress entities.EmailAddress) (bool, error) {
	person, err := ctx.PersonRepository().GetByEmail(emailAddress)
	if err != nil {
		return false, err
	}
	return person == nil, nil
}
