package attendance

import (
	"context"

	"github.com/Goldwin/ies-pik-cms/pkg/attendance/commands"
	"github.com/Goldwin/ies-pik-cms/pkg/attendance/dto"
	"github.com/Goldwin/ies-pik-cms/pkg/attendance/entities"
	"github.com/Goldwin/ies-pik-cms/pkg/attendance/queries"
	"github.com/Goldwin/ies-pik-cms/pkg/common/out"
	"github.com/Goldwin/ies-pik-cms/pkg/common/utils"
	"github.com/Goldwin/ies-pik-cms/pkg/common/worker"
	"github.com/samber/lo"
)

type AttendanceDataLayerComponent interface {
	QueryWorker() worker.QueryWorker[queries.QueryContext]
	CommandWorker() worker.UnitOfWork[commands.CommandContext]
}

type AttendanceCommandComponent interface {
	CreateEventSchedule(ctx context.Context, schedule dto.EventScheduleDTO, output out.Output[dto.EventScheduleDTO]) out.Waitable
	AddEventScheduleActivity(ctx context.Context, activity dto.EventScheduleActivityDTO, output out.Output[dto.EventScheduleDTO]) out.Waitable
	CreateNextEvent(ctx context.Context, scheduleID string, output out.Output[[]dto.EventDTO]) out.Waitable
	UpdateEventSchedule(ctx context.Context, schedule dto.EventScheduleDTO, output out.Output[dto.EventScheduleDTO]) out.Waitable
	RemoveEventScheduleActivity(ctx context.Context, activity dto.EventScheduleActivityDTO, output out.Output[dto.EventScheduleDTO]) out.Waitable
	UpdateEventScheduleActivity(ctx context.Context, activity dto.EventScheduleActivityDTO, output out.Output[dto.EventScheduleDTO]) out.Waitable
	HouseholdCheckin(ctx context.Context, checkins dto.HouseholdCheckinDTO, output out.Output[[]dto.EventAttendanceDTO]) out.Waitable
	CheckIn(ctx context.Context, attendance dto.EventAttendanceDTO, output out.Output[dto.EventAttendanceDTO]) out.Waitable
}

type AttendanceQueryComponent interface {
	GetEventSchedule(ctx context.Context, query queries.GetEventScheduleFilter, output out.Output[queries.GetEventScheduleResult])
	ListEventSchedules(ctx context.Context, query queries.ListEventScheduleFilter, output out.Output[queries.ListEventScheduleResult])
	ListEventsBySchedule(ctx context.Context, query queries.ListEventByScheduleFilter, output out.Output[queries.ListEventByScheduleResult])
	GetEvent(ctx context.Context, query queries.GetEventFilter, output out.Output[queries.GetEventResult])
	ListEventAttendance(ctx context.Context, query queries.ListEventAttendanceFilter, output out.Output[queries.ListEventAttendanceResult])
	SearchHousehold(ctx context.Context, filter queries.SearchHouseholdFilter, output out.Output[queries.SearchHouseholdResult])
}

type AttendanceComponent interface {
	AttendanceCommandComponent
	AttendanceQueryComponent
}

type attendanceComponentImpl struct {
	dataLayer AttendanceDataLayerComponent
}

// SearchHousehold implements AttendanceComponent.
func (a *attendanceComponentImpl) SearchHousehold(ctx context.Context, filter queries.SearchHouseholdFilter, output out.Output[queries.SearchHouseholdResult]) {
	result, err := a.dataLayer.QueryWorker().Query(ctx).SearchHousehold().Execute(filter)

	if err.NoError() {
		output.OnSuccess(result)
		return
	}
	output.OnError(out.ConvertQueryErrorDetail(err))
}

// HouseholdCheckin implements AttendanceComponent.
func (a *attendanceComponentImpl) HouseholdCheckin(ctx context.Context, input dto.HouseholdCheckinDTO, output out.Output[[]dto.EventAttendanceDTO]) out.Waitable {
	x := commands.HouseholdCheckinCommand{
		EventID: input.EventID,
		Attendee: lo.Map(input.Attendees, func(e dto.PersonCheckinDTO, _ int) commands.Attendee {
			return commands.Attendee{
				PersonID:       e.PersonID,
				ActivityID:     e.EventActivityID,
				AttendanceType: e.AttendanceType,
			}
		}),
	}
	return utils.SingleCommandExecution(a.dataLayer.CommandWorker(), x).WithOutput(
		out.OutputAdapter(output, func(e []*entities.Attendance) []dto.EventAttendanceDTO {
			return lo.Map(e, func(e *entities.Attendance, _ int) dto.EventAttendanceDTO {
				return dto.FromAttendanceEntities(e)
			})
		}),
	).Execute(ctx)
}

// GetEvent implements AttendanceComponent.
func (a *attendanceComponentImpl) GetEvent(ctx context.Context, query queries.GetEventFilter, output out.Output[queries.GetEventResult]) {

	result, err := a.dataLayer.QueryWorker().Query(ctx).GetEvent().Execute(query)
	if err.NoError() {
		output.OnSuccess(result)
		return
	}
	output.OnError(out.ConvertQueryErrorDetail(err))

}

// GetEventSchedule implements AttendanceComponent.
func (a *attendanceComponentImpl) GetEventSchedule(ctx context.Context, query queries.GetEventScheduleFilter, output out.Output[queries.GetEventScheduleResult]) {

	result, err := a.dataLayer.QueryWorker().Query(ctx).GetEventSchedule().Execute(query)
	if err.NoError() {
		output.OnSuccess(result)
		return
	}
	output.OnError(out.ConvertQueryErrorDetail(err))

}

// ListEventAttendance implements AttendanceComponent.
func (a *attendanceComponentImpl) ListEventAttendance(ctx context.Context, query queries.ListEventAttendanceFilter, output out.Output[queries.ListEventAttendanceResult]) {
	result, err := a.dataLayer.QueryWorker().Query(ctx).ListEventAttendance().Execute(query)

	if err.NoError() {
		output.OnSuccess(result)
		return
	}
	output.OnError(out.ConvertQueryErrorDetail(err))

}

// ListEventSchedules implements AttendanceComponent.
func (a *attendanceComponentImpl) ListEventSchedules(ctx context.Context, query queries.ListEventScheduleFilter, output out.Output[queries.ListEventScheduleResult]) {

	result, err := a.dataLayer.QueryWorker().Query(ctx).ListEventSchedules().Execute(queries.ListEventScheduleFilter{Limit: query.Limit, LastID: query.LastID})
	if err.NoError() {
		output.OnSuccess(result)
		return
	}

	output.OnError(out.ConvertQueryErrorDetail(err))

}

// ListEventsBySchedule implements AttendanceComponent.
func (a *attendanceComponentImpl) ListEventsBySchedule(ctx context.Context, query queries.ListEventByScheduleFilter, output out.Output[queries.ListEventByScheduleResult]) {
	result, err := a.dataLayer.QueryWorker().Query(ctx).ListEventsBySchedule().Execute(query)
	if err.NoError() {
		output.OnSuccess(result)
		return
	}
	output.OnError(out.ConvertQueryErrorDetail(err))

}

// UpdateEventSchedule implements AttendanceComponent.
func (a *attendanceComponentImpl) UpdateEventSchedule(ctx context.Context, schedule dto.EventScheduleDTO, output out.Output[dto.EventScheduleDTO]) out.Waitable {
	return utils.SingleCommandExecution(a.dataLayer.CommandWorker(), commands.UpdateEventScheduleCommand{
		ID:             schedule.ID,
		Name:           schedule.Name,
		ScheduleType:   string(schedule.Type),
		TimezoneOffset: schedule.TimezoneOffset,
		Days:           schedule.Days,
		Date:           schedule.Date,
		StartDate:      schedule.StartDate,
		EndDate:        schedule.EndDate,
	}).WithOutput(
		out.OutputAdapter(output, func(e entities.EventSchedule) dto.EventScheduleDTO {
			return dto.FromEntities(&e)
		}),
	).Execute(ctx)
}

// RemoveEventScheduleActivity implements AttendanceComponent.
func (a *attendanceComponentImpl) RemoveEventScheduleActivity(ctx context.Context, activity dto.EventScheduleActivityDTO, output out.Output[dto.EventScheduleDTO]) out.Waitable {
	return utils.SingleCommandExecution(a.dataLayer.CommandWorker(), commands.RemoveScheduleActivityCommand{
		ScheduleID: activity.ScheduleID,
		ActivityID: activity.ID,
	}).WithOutput(
		out.OutputAdapter(output, func(e entities.EventSchedule) dto.EventScheduleDTO {
			return dto.FromEntities(&e)
		}),
	).Execute(ctx)
}

// AddEventScheduleActivity implements AttendanceComponent.
func (a *attendanceComponentImpl) AddEventScheduleActivity(ctx context.Context, activity dto.EventScheduleActivityDTO, output out.Output[dto.EventScheduleDTO]) out.Waitable {
	return utils.SingleCommandExecution(a.dataLayer.CommandWorker(), commands.AddEventScheduleActivityCommand{
		ScheduleID: activity.ScheduleID,
		Name:       activity.Name,
		Hour:       activity.Hour,
		Minute:     activity.Minute,
	}).WithOutput(
		out.OutputAdapter(output, func(e entities.EventSchedule) dto.EventScheduleDTO {
			return dto.FromEntities(&e)
		}),
	).Execute(ctx)
}

// CreateNextEvent implements AttendanceComponent.
func (a *attendanceComponentImpl) CreateNextEvent(ctx context.Context, scheduleID string, output out.Output[[]dto.EventDTO]) out.Waitable {
	return utils.SingleCommandExecution(a.dataLayer.CommandWorker(), commands.CreateNextEventCommand{
		ScheduleID: scheduleID,
	}).WithOutput(
		out.OutputAdapter(output, func(e []*entities.Event) []dto.EventDTO {
			return lo.Map(e, func(e *entities.Event, _ int) dto.EventDTO {
				return dto.FromEventEntities(e)
			})
		}),
	).Execute(ctx)
}

// CreateEventSchedule implements AttendanceComponent.
func (a *attendanceComponentImpl) CreateEventSchedule(ctx context.Context, schedule dto.EventScheduleDTO, output out.Output[dto.EventScheduleDTO]) out.Waitable {
	return utils.SingleCommandExecution(a.dataLayer.CommandWorker(), commands.CreateEventScheduleCommand{
		Name:           schedule.Name,
		ScheduleType:   schedule.Type,
		TimezoneOffset: schedule.TimezoneOffset,
		Days:           schedule.Days,
		Date:           schedule.Date,
		StartDate:      schedule.StartDate,
		EndDate:        schedule.EndDate,
	}).WithOutput(
		out.OutputAdapter(output, func(e entities.EventSchedule) dto.EventScheduleDTO {
			return dto.FromEntities(&e)
		}),
	).Execute(ctx)
}

// UpdateEventScheduleActivity implements AttendanceComponent.
func (a *attendanceComponentImpl) UpdateEventScheduleActivity(ctx context.Context, activity dto.EventScheduleActivityDTO, output out.Output[dto.EventScheduleDTO]) out.Waitable {
	return utils.SingleCommandExecution(a.dataLayer.CommandWorker(), commands.UpdateEventScheduleActivityCommand{
		ScheduleID: activity.ScheduleID,
		ActivityID: activity.ID,
		Name:       activity.Name,
		Hour:       activity.Hour,
		Minute:     activity.Minute,
	}).WithOutput(
		out.OutputAdapter(output, func(e entities.EventSchedule) dto.EventScheduleDTO {
			return dto.FromEntities(&e)
		})).Execute(ctx)
}

// CheckIn implements AttendanceComponent.
func (a *attendanceComponentImpl) CheckIn(ctx context.Context, attendance dto.EventAttendanceDTO, output out.Output[dto.EventAttendanceDTO]) out.Waitable {
	return utils.SingleCommandExecution(a.dataLayer.CommandWorker(), commands.CheckInCommand{
		EventID:    attendance.Event.ID,
		ActivityID: attendance.Activity.ID,
		Person: commands.PersonInput{
			PersonID:          attendance.Person.PersonID,
			FirstName:         attendance.Person.FirstName,
			MiddleName:        attendance.Person.MiddleName,
			LastName:          attendance.Person.LastName,
			ProfilePictureUrl: attendance.Person.ProfilePictureURL,
		},
		Type: attendance.AttendanceType,
	}).WithOutput(
		out.OutputAdapter(output, func(e entities.Attendance) dto.EventAttendanceDTO {
			return dto.FromAttendanceEntities(&e)
		})).Execute(ctx)
}

func NewAttendanceComponent(datalayer AttendanceDataLayerComponent) AttendanceComponent {
	return &attendanceComponentImpl{
		dataLayer: datalayer,
	}
}
