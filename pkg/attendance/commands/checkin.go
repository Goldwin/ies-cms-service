package commands

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/Goldwin/ies-pik-cms/pkg/attendance/entities"
	. "github.com/Goldwin/ies-pik-cms/pkg/common/commands"
	"github.com/samber/lo"
)

const (
	CheckInEventDoesNotExistsError    CommandErrorCode = 30401
	CheckInActivityDoesNotExistsError CommandErrorCode = 30402
	CheckInInvalidInputError          CommandErrorCode = 30403
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
	Type       string
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

	activity, found := lo.Find(event.EventActivities, func(e *entities.EventActivity) bool {
		return e.ID == c.ActivityID
	})

	if !found {
		return CommandExecutionResult[entities.Attendance]{
			Status: ExecutionStatusFailed,
			Error:  CommandErrorDetail{Code: CheckInActivityDoesNotExistsError, Message: "Activity not found"},
		}
	}

	person, err := ctx.PersonRepository().Get(c.Person.PersonID)

	if err != nil {
		return CommandExecutionResult[entities.Attendance]{
			Status: ExecutionStatusFailed,
			Error:  CommandErrorDetailWorkerFailure(err),
		}
	}

	if person == nil {
		return CommandExecutionResult[entities.Attendance]{
			Status: ExecutionStatusFailed,
			Error:  CommandErrorDetail{Code: CheckInInvalidInputError, Message: "Person not found"},
		}
	}

	securityCode := lo.RandomString(5, []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"))
	securityNumber := rand.Int() % 1000

	attendanceID := fmt.Sprintf("%s.%s", event.ID, c.Person.PersonID)
	attendance := &entities.Attendance{
		ID:             attendanceID,
		Event:          event,
		EventActivity:  activity,
		SecurityCode:   securityCode,
		SecurityNumber: securityNumber,
		CheckinTime:    time.Now(),
		CheckoutTime:   time.UnixMilli(0),
		Type:           entities.AttendanceType(c.Type),
	}

	validationErr := attendance.IsValid()

	if validationErr != "" {
		return CommandExecutionResult[entities.Attendance]{
			Status: ExecutionStatusFailed,
			Error:  CommandErrorDetail{Code: CheckInInvalidInputError, Message: validationErr},
		}
	}

	attendance, err = ctx.AttendanceRepository().Save(attendance)

	if err != nil {
		return CommandExecutionResult[entities.Attendance]{
			Status: ExecutionStatusFailed,
			Error:  CommandErrorDetailWorkerFailure(err),
		}
	}

	return CommandExecutionResult[entities.Attendance]{
		Status: ExecutionStatusSuccess,
		Result: *attendance,
	}

}
