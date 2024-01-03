package commands

import (
	"fmt"

	. "github.com/Goldwin/ies-pik-cms/pkg/common/commands"
	"github.com/Goldwin/ies-pik-cms/pkg/people/dto"
	"github.com/Goldwin/ies-pik-cms/pkg/people/entities"
)

const (
	UpdateHouseholdErrorCodeVerifyDataError                  AppErrorCode = 10201
	UpdateHouseholdErrorCodePersonNotExistsError             AppErrorCode = 10202
	UpdateHouseholdErrorCodeOneOrMorePersonHasHouseholdError AppErrorCode = 10203
)

type UpdateHouseholdCommand struct {
	Input dto.HouseHoldInput
}

func (cmd UpdateHouseholdCommand) Execute(ctx CommandContext) AppExecutionResult[dto.Household] {
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
		ID:            cmd.Input.ID,
	}

	result, err := ctx.HouseholdRepository().UpdateHousehold(household)

	memberPhone := ""
	if len(result.HouseholdHead.PhoneNumbers) > 0 {
		memberPhone = string(result.HouseholdHead.PhoneNumbers[0])
	}

	householdHeadDto := dto.HouseholdPerson{
		ID:           result.HouseholdHead.ID,
		FirstName:    result.HouseholdHead.FirstName,
		MiddleName:   result.HouseholdHead.MiddleName,
		LastName:     result.HouseholdHead.LastName,
		EmailAddress: string(result.HouseholdHead.EmailAddress),
		PhoneNumber:  string(memberPhone),
	}

	householdMembersDto := make([]dto.HouseholdPerson, len(result.Members))

	for i, member := range result.Members {

		memberPhone = ""
		if len(result.HouseholdHead.PhoneNumbers) > 0 {
			memberPhone = string(result.HouseholdHead.PhoneNumbers[0])
		}
		householdMembersDto[i] = dto.HouseholdPerson{
			ID:           member.ID,
			FirstName:    member.FirstName,
			LastName:     member.LastName,
			EmailAddress: string(member.EmailAddress),
			PhoneNumber:  string(memberPhone),
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
