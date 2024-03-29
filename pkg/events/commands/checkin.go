package commands

import (
	"fmt"
	"time"

	. "github.com/Goldwin/ies-pik-cms/pkg/common/commands"
	"github.com/Goldwin/ies-pik-cms/pkg/events/dto"
	"github.com/Goldwin/ies-pik-cms/pkg/events/entities"
)

const (
	CheckInCommandsFailedToFetchPerson CommandErrorCode = 30001
	CheckInCommandsErrorPersonNotFound CommandErrorCode = 30002
	CheckInCommandsFailedToFetchEvent  CommandErrorCode = 30003
	CheckInCommandsErrorEventNotFound  CommandErrorCode = 30004
	CheckInCommandsFailedToCheckIn     CommandErrorCode = 30005
)

type CheckInCommands struct {
	Input dto.CheckInInput
}

func (cmd CheckInCommands) Execute(ctx CommandContext) CommandExecutionResult[[]dto.CheckInEvent] {
	persons, err := ctx.PersonRepository().GetByIds(cmd.Input.PersonIDs)
	if err != nil {
		return CommandExecutionResult[[]dto.CheckInEvent]{
			Status: ExecutionStatusFailed,
			Error: CommandErrorDetail{
				Code:    CheckInCommandsFailedToFetchPerson,
				Message: fmt.Sprintf("Failed to Fetch Person: %s", err.Error()),
			},
		}
	}

	if len(persons) != len(cmd.Input.PersonIDs) {
		return CommandExecutionResult[[]dto.CheckInEvent]{
			Status: ExecutionStatusFailed,
			Error: CommandErrorDetail{
				Code:    CheckInCommandsErrorPersonNotFound,
				Message: fmt.Sprintf("Cannot fetch all person info"),
			},
		}
	}

	checker, err := ctx.PersonRepository().Get(cmd.Input.CheckerID)

	if err != nil {
		return CommandExecutionResult[[]dto.CheckInEvent]{
			Status: ExecutionStatusFailed,
			Error: CommandErrorDetail{
				Code:    CheckInCommandsFailedToFetchPerson,
				Message: fmt.Sprintf("Failed to Fetch Person: %s", err.Error()),
			},
		}
	}

	if checker == nil {
		return CommandExecutionResult[[]dto.CheckInEvent]{
			Status: ExecutionStatusFailed,
			Error: CommandErrorDetail{
				Code:    CheckInCommandsErrorPersonNotFound,
				Message: fmt.Sprintf("Person Not Found"),
			},
		}
	}

	session, err := ctx.ChurchEventSessionRepository().Get(cmd.Input.EventID)

	if err != nil {
		return CommandExecutionResult[[]dto.CheckInEvent]{
			Status: ExecutionStatusFailed,
			Error: CommandErrorDetail{
				Code:    CheckInCommandsFailedToFetchEvent,
				Message: fmt.Sprintf("Failed to Fetch Event: %s", err.Error()),
			},
		}
	}

	if session == nil {
		return CommandExecutionResult[[]dto.CheckInEvent]{
			Status: ExecutionStatusFailed,
			Error: CommandErrorDetail{
				Code:    CheckInCommandsErrorEventNotFound,
				Message: fmt.Sprintf("Event Not Found: %s", err.Error()),
			},
		}
	}

	result := []dto.CheckInEvent{}

	for _, person := range persons {
		checkin := entities.CheckInEvent{
			ID:        fmt.Sprintf("%s.%s", session.ID, person.ID),
			Person:    person,
			SessionID: session.ID,
			CheckInAt: time.Now(),
		}
		err = ctx.EventCheckInRepository().Save(checkin)
		if err != nil {
			return CommandExecutionResult[[]dto.CheckInEvent]{
				Status: ExecutionStatusFailed,
				Error: CommandErrorDetail{
					Code:    CheckInCommandsFailedToCheckIn,
					Message: fmt.Sprintf("Failed to CheckIn: %s", err.Error()),
				},
			}
		}
		result = append(result, dto.CheckInEvent{
			ID: checkin.SessionID,
			Person: dto.Person{
				ID:        checkin.Person.ID,
				FirstName: checkin.Person.FirstName,
				LastName:  checkin.Person.LastName,
			},
			SessionID: checkin.SessionID,
			CheckInAt: checkin.CheckInAt,
		})
	}

	return CommandExecutionResult[[]dto.CheckInEvent]{
		Status: ExecutionStatusSuccess,
		Result: result,
	}
}
