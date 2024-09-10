package commands

import (
	"math/rand"
	"time"

	"github.com/Goldwin/ies-pik-cms/pkg/attendance/entities"
	. "github.com/Goldwin/ies-pik-cms/pkg/common/commands"
	"github.com/samber/lo"
)

const (
	HouseholdCheckinCommandEventDoesNotExistsError = 30501
	HouseholdCheckinCommandPersonMissingError      = 30501
)

var (
	charset = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
)

type Attendee struct {
	PersonID       string
	ActivityID     string
	AttendanceType string
}

type HouseholdCheckinCommand struct {
	EventID     string
	Attendee    []Attendee
	CheckedInBy string
}

func (c HouseholdCheckinCommand) Execute(ctx CommandContext) CommandExecutionResult[[]*entities.Attendance] {
	event, err := ctx.EventRepository().Get(c.EventID)

	if err != nil {
		return CommandExecutionResult[[]*entities.Attendance]{
			Status: ExecutionStatusFailed,
			Error:  CommandErrorDetailWorkerFailure(err),
		}
	}

	if event == nil {
		return CommandExecutionResult[[]*entities.Attendance]{
			Status: ExecutionStatusFailed,
			Error:  CommandErrorDetail{Code: HouseholdCheckinCommandEventDoesNotExistsError, Message: "Event not found"},
		}
	}

	attendees, err := ctx.PersonRepository().List(lo.Map(c.Attendee, func(e Attendee, _ int) string { return e.PersonID }))
	if err != nil {
		return CommandExecutionResult[[]*entities.Attendance]{
			Status: ExecutionStatusFailed,
			Error:  CommandErrorDetailWorkerFailure(err),
		}
	}

	if len(attendees) != len(c.Attendee) {
		return CommandExecutionResult[[]*entities.Attendance]{
			Status: ExecutionStatusFailed,
			Error:  CommandErrorDetail{Code: HouseholdCheckinCommandEventDoesNotExistsError, Message: "One or more Person not found"},
		}
	}

	checkinPerson, err := ctx.PersonRepository().Get(c.CheckedInBy)
	if err != nil {
		return CommandExecutionResult[[]*entities.Attendance]{
			Status: ExecutionStatusFailed,
			Error:  CommandErrorDetailWorkerFailure(err),
		}
	}

	if checkinPerson == nil {
		return CommandExecutionResult[[]*entities.Attendance]{
			Status: ExecutionStatusFailed,
			Error:  CommandErrorDetail{Code: HouseholdCheckinCommandPersonMissingError, Message: "Can't Find person who check-in"},
		}
	}

	activitiesMap := lo.SliceToMap(event.EventActivities, func(e *entities.EventActivity) (string, *entities.EventActivity) { return e.ID, e })

	attendeesMap := lo.SliceToMap(attendees, func(e *entities.Person) (string, *entities.Person) { return e.PersonID, e })

	attendances := lo.Map(c.Attendee, func(a Attendee, _ int) *entities.Attendance {
		securityCode := lo.RandomString(5, charset)
		securityNumber := rand.Int() % 1000
		return &entities.Attendance{
			ID:             c.EventID + "." + a.ActivityID + "." + a.PersonID,
			Event:          event,
			EventActivity:  activitiesMap[a.ActivityID],
			Attendee:       attendeesMap[a.PersonID],
			CheckedInBy:    checkinPerson,
			SecurityCode:   securityCode,
			SecurityNumber: securityNumber,
			CheckinTime:    time.Now(),
			Type:           entities.AttendanceType(a.AttendanceType),
		}
	})

	for _, attendance := range attendances {
		_, err = ctx.AttendanceRepository().Save(attendance)
		if err != nil {
			return CommandExecutionResult[[]*entities.Attendance]{
				Status: ExecutionStatusFailed,
				Error:  CommandErrorDetailWorkerFailure(err),
			}
		}
	}

	return CommandExecutionResult[[]*entities.Attendance]{
		Status: ExecutionStatusSuccess,
		Result: attendances,
	}
}
