package commands

import (
	"fmt"

	. "github.com/Goldwin/ies-pik-cms/pkg/common/commands"
	"github.com/Goldwin/ies-pik-cms/pkg/people/dto"
	"github.com/Goldwin/ies-pik-cms/pkg/people/entities"
)

const (
	AddHouseholdErrorCodeDBError                          AppErrorCode = 10201
	AddHouseholdErrorCodePersonNotExistsError             AppErrorCode = 10202
	AddHouseholdErrorCodeOneOrMorePersonHasHouseholdError AppErrorCode = 10203
)

type AddHouseholdCommand struct {
	Input dto.HouseHoldInput
}

func (cmd AddHouseholdCommand) Execute(ctx CommandContext) AppExecutionResult[dto.Household] {
	householdHead, err := ctx.PersonRepository().Get(cmd.Input.HouseholdHeadPersonId)

	if err != nil {
		return AppExecutionResult[dto.Household]{
			Status: ExecutionStatusFailed,
			Error: AppErrorDetail{
				Code:    AddHouseholdErrorCodeDBError,
				Message: fmt.Sprintf("Can't Add New Household Info, Error: %s", err.Error()),
			},
		}
	}

	if householdHead == nil {
		return AppExecutionResult[dto.Household]{
			Status: ExecutionStatusFailed,
			Error: AppErrorDetail{
				Code:    AddHouseholdErrorCodeDBError,
				Message: fmt.Sprintf("Can't Add New Household Info, Error: Person with id %s does not exist", cmd.Input.HouseholdHeadPersonId),
			},
		}
	}

	persons, err := ctx.PersonRepository().ListByID(cmd.Input.MemberPersonsIds)
	if err != nil {
		return AppExecutionResult[dto.Household]{
			Status: ExecutionStatusFailed,
			Error: AppErrorDetail{
				Code:    AddHouseholdErrorCodeDBError,
				Message: fmt.Sprintf("Can't Add New Household Info, Error: %s", err.Error()),
			},
		}
	}

	personIdMap := map[string]bool{}
	for _, personId := range cmd.Input.MemberPersonsIds {
		personIdMap[personId] = true
	}

	for _, person := range persons {
		if _, ok := personIdMap[person.ID]; !ok {
			return AppExecutionResult[dto.Household]{
				Status: ExecutionStatusFailed,
				Error: AppErrorDetail{
					Code:    AddHouseholdErrorCodeDBError,
					Message: fmt.Sprintf("Can't Add New Household Info, Error: Person with id %s does not exist", person.ID),
				},
			}
		}
	}

	household := entities.Household{
		HouseholdHead: *householdHead,
		Members:       persons,
		Name:          cmd.Input.Name,
		PictureUrl:    cmd.Input.PictureUrl,
	}

	result, err := ctx.HouseholdRepository().AddHousehold(household)

	householdHeadDto := dto.HouseholdPerson{
		ID:           result.HouseholdHead.ID,
		FirstName:    result.HouseholdHead.FirstName,
		MiddleName:   result.HouseholdHead.MiddleName,
		LastName:     result.HouseholdHead.LastName,
		EmailAddress: string(result.HouseholdHead.EmailAddress),
		PhoneNumber:  string(result.HouseholdHead.PhoneNumber),
	}

	householdMembersDto := make([]dto.HouseholdPerson, len(result.Members))

	for i, member := range result.Members {

		householdMembersDto[i] = dto.HouseholdPerson{
			ID:           member.ID,
			FirstName:    member.FirstName,
			LastName:     member.LastName,
			EmailAddress: string(member.EmailAddress),
			PhoneNumber:  string(member.PhoneNumber),
		}
	}

	return AppExecutionResult[dto.Household]{
		Status: ExecutionStatusSuccess,
		Error:  AppErrorDetailNone,
		Result: dto.Household{
			ID:            result.ID,
			HouseholdHead: householdHeadDto,
			Members:       householdMembersDto,
			PictureUrl:    result.PictureUrl,
			Name:          result.Name,
		},
	}
}
