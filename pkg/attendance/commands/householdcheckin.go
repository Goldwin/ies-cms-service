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
	HouseholdCheckinCommandEventDoesNotExistsError = 30501
	HouseholdCheckinCommandPersonMissingError      = 30501
	HouseholdCheckinCommandEventEndedError         = 30503
	HouseholdCheckinCommandInvalidAttendanceError  = 30504
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

type invalidAttendance struct {
	person Attendee
	reason string
}

func (c HouseholdCheckinCommand) Execute(ctx CommandContext) CommandExecutionResult[[]*entities.Attendance] {
	checkinTime := time.Now()
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

	if event.EndDate.Before(checkinTime) {
		return CommandExecutionResult[[]*entities.Attendance]{
			Status: ExecutionStatusFailed,
			Error:  CommandErrorDetail{Code: HouseholdCheckinCommandEventEndedError, Message: "Failed to checkin. Event already ended"},
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

	personAttendanceSummaries, err := ctx.PersonAttendanceSummaryRepository().List(
		lo.Map(attendees, func(e *entities.Person, _ int) string {
			return entities.PersonAttendanceSummary{PersonID: e.PersonID, ScheduleID: event.ScheduleID}.ID()
		}),
	)

	if err != nil {
		return CommandExecutionResult[[]*entities.Attendance]{
			Status: ExecutionStatusFailed,
			Error:  CommandErrorDetailWorkerFailure(err),
		}
	}

	personToAttendanceSummaryMap := lo.SliceToMap(personAttendanceSummaries, func(e *entities.PersonAttendanceSummary) (string, *entities.PersonAttendanceSummary) {
		return e.PersonID, e
	})

	activitiesMap := lo.SliceToMap(event.EventActivities, func(e *entities.EventActivity) (string, *entities.EventActivity) { return e.ID, e })

	attendeesMap := lo.SliceToMap(attendees, func(e *entities.Person) (string, *entities.Person) { return e.PersonID, e })

	invalidAttendances := []invalidAttendance{}

	attendances := lo.Map(c.Attendee, func(a Attendee, _ int) *entities.Attendance {
		securityCode := lo.RandomString(5, charset)
		securityNumber := rand.Int() % 1000
		activity, ok := activitiesMap[a.ActivityID]
		if !ok {
			invalidAttendances = append(invalidAttendances, invalidAttendance{
				person: a,
				reason: fmt.Sprintf("Activity %s not found", a.ActivityID),
			})
		}
		attendee, ok := attendeesMap[a.PersonID]
		if !ok {
			invalidAttendances = append(invalidAttendances, invalidAttendance{
				person: a,
				reason: fmt.Sprintf("Person %s not found", a.PersonID),
			})
		}
		return &entities.Attendance{
			ID:             c.EventID + "." + a.ActivityID + "." + a.PersonID,
			Event:          event,
			EventActivity:  activity,
			Attendee:       attendee,
			CheckedInBy:    checkinPerson,
			SecurityCode:   securityCode,
			SecurityNumber: securityNumber,
			CheckinTime:    checkinTime,
			Type:           entities.AttendanceType(a.AttendanceType),
			FirstTime:      false,
		}
	})

	if len(invalidAttendances) > 0 {
		return CommandExecutionResult[[]*entities.Attendance]{
			Status: ExecutionStatusFailed,
			Error:  CommandErrorDetail{Code: HouseholdCheckinCommandInvalidAttendanceError, Message: "Invalid Attendances", Details: lo.Map(invalidAttendances, func(e invalidAttendance, _ int) string { return e.reason })},
		}
	}

	for _, attendance := range attendances {
		summary, ok := personToAttendanceSummaryMap[attendance.Attendee.PersonID]
		if !ok {
			summary = &entities.PersonAttendanceSummary{
				PersonID:             attendance.Attendee.PersonID,
				ScheduleID:           event.ScheduleID,
				FirstEventAttendance: attendance,
			}
			attendance.FirstTime = true
		}
		summary.LastEventAttendance = attendance
		_, err = ctx.AttendanceRepository().Save(attendance)
		if err != nil {
			return CommandExecutionResult[[]*entities.Attendance]{
				Status: ExecutionStatusFailed,
				Error:  CommandErrorDetailWorkerFailure(err),
			}
		}
		_, err = ctx.PersonAttendanceSummaryRepository().Save(summary)

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
