package commands

import (
	"time"

	"github.com/Goldwin/ies-pik-cms/pkg/attendance/entities"
	. "github.com/Goldwin/ies-pik-cms/pkg/common/commands"
	"github.com/samber/lo"
)

const (
	CheckInEventDoesNotExistsError    CommandErrorCode = 30401
	CheckInActivityDoesNotExistsError CommandErrorCode = 30402
)

type PersonInput struct {
	PersonID          string
	FirstName         string
	MiddleName        string
	LastName          string
	ProfilePictureUrl string
}

type CheckInCommand struct {
	EventID    string
	ActivityID string
	Person     PersonInput
}

func (c CheckInCommand) Execute(ctx CommandContext) CommandExecutionResult[entities.Attendance] {

	event, err := ctx.EventRepository().Get(c.EventID)

	if err != nil {
		return CommandExecutionResult[entities.Attendance]{
			Status: ExecutionStatusFailed,
			Error:  CommandErrorDetailWorkerFailure(err),
		}
	}

	if event == nil {
		return CommandExecutionResult[entities.Attendance]{
			Status: ExecutionStatusFailed,
			Error:  CommandErrorDetail{Code: CheckInEventDoesNotExistsError, Message: "Schedule not found"},
		}
	}

	activity, found := lo.Find(event.EventActivities, func(e entities.EventActivity) bool {
		return e.ID == c.ActivityID
	})

	if !found {
		return CommandExecutionResult[entities.Attendance]{
			Status: ExecutionStatusFailed,
			Error:  CommandErrorDetail{Code: CheckInActivityDoesNotExistsError, Message: "Activity not found"},
		}
	}

	//TODO upsert attendance.
	return CommandExecutionResult[entities.Attendance]{
		Status: ExecutionStatusSuccess,
		Result: entities.Attendance{
			ID:                "",
			EventID:           event.ID,
			EventActivityID:   activity.ID,
			PersonID:          "",
			FirstName:         "",
			MiddleName:        "",
			LastName:          "",
			ProfilePictureUrl: "",
			SecurityCode:      "",
			SecurityNumber:    0,
			CheckinTime:       time.Now(),
			CheckoutTime:      time.Now(),
			Type:              "",
		},
	}

}
