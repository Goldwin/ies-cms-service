package commands

import (
	"fmt"
	"time"

	. "github.com/Goldwin/ies-pik-cms/pkg/common/commands"
	"github.com/Goldwin/ies-pik-cms/pkg/events/dto"
	"github.com/Goldwin/ies-pik-cms/pkg/events/entities"
	"github.com/Goldwin/ies-pik-cms/pkg/events/repositories"
)

const (
	CheckInCommandsFailedToFetchPerson AppErrorCode = 30001
	CheckInCommandsErrorPersonNotFound AppErrorCode = 30002
	CheckInCommandsFailedToFetchEvent  AppErrorCode = 30003
	CheckInCommandsErrorEventNotFound  AppErrorCode = 30004
	CheckInCommandsFailedToCheckIn     AppErrorCode = 30005
)

type CheckInCommands struct {
	Input dto.CheckInInput
}

func (cmd CheckInCommands) Execute(ctx repositories.CommandContext) AppExecutionResult[dto.CheckInEvent] {
	person, err := ctx.PersonRepository().Get(cmd.Input.PersonID)
	if err != nil {
		return AppExecutionResult[dto.CheckInEvent]{
			Status: ExecutionStatusFailed,
			Error: AppErrorDetail{
				Code:    CheckInCommandsFailedToFetchPerson,
				Message: fmt.Sprintf("Failed to Fetch Person: %s", err.Error()),
			},
		}
	}

	if person == nil {
		return AppExecutionResult[dto.CheckInEvent]{
			Status: ExecutionStatusFailed,
			Error: AppErrorDetail{
				Code:    CheckInCommandsErrorPersonNotFound,
				Message: fmt.Sprintf("Person Not Found: %s", err.Error()),
			},
		}
	}

	event, err := ctx.ChurchEventRepository().Get(cmd.Input.EventID)

	if err != nil {
		return AppExecutionResult[dto.CheckInEvent]{
			Status: ExecutionStatusFailed,
			Error: AppErrorDetail{
				Code:    CheckInCommandsFailedToFetchEvent,
				Message: fmt.Sprintf("Failed to Fetch Event: %s", err.Error()),
			},
		}
	}

	if event == nil {
		return AppExecutionResult[dto.CheckInEvent]{
			Status: ExecutionStatusFailed,
			Error: AppErrorDetail{
				Code:    CheckInCommandsErrorEventNotFound,
				Message: fmt.Sprintf("Event Not Found: %s", err.Error()),
			},
		}
	}

	checkin := entities.CheckInEvent{
		ID:        fmt.Sprintf("%s.%s", event.ID, person.ID),
		Person:    *person,
		Event:     *event,
		CheckInAt: time.Now(),
	}

	err = ctx.EventCheckInRepository().Save(checkin)

	if err != nil {
		return AppExecutionResult[dto.CheckInEvent]{
			Status: ExecutionStatusFailed,
			Error: AppErrorDetail{
				Code:    CheckInCommandsFailedToCheckIn,
				Message: fmt.Sprintf("Failed to CheckIn: %s", err.Error()),
			},
		}
	}

	return AppExecutionResult[dto.CheckInEvent]{
		Status: ExecutionStatusSuccess,
		Result: dto.CheckInEvent{
			ID: checkin.ID,
			Person: dto.Person{
				ID:        person.ID,
				FirstName: person.FirstName,
				LastName:  person.LastName,
			},
			Event: dto.ChurchEvent{
				ID:        event.ID,
				Name:      event.Name,
				StartTime: event.StartTime,
			},
			CheckInAt: checkin.CheckInAt,
		},
	}
}
