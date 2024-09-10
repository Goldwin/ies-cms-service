package mongo

import (
	"time"

	"github.com/Goldwin/ies-pik-cms/pkg/attendance/dto"
	"github.com/Goldwin/ies-pik-cms/pkg/attendance/entities"
	"github.com/samber/lo"
)

const (
	AttendanceCollection        = "attendances"
	AttendanceSummaryCollection = "attendance_summaries"
	EventScheduleCollection     = "event_schedules"
	EventCollection             = "events"
	PersonCollection            = "persons"
	HouseholdCollection         = "households"
	PersonHouseholdCollection   = "person_households"
)

type EventScheduleModel struct {
	ID             string                       `bson:"_id"`
	Name           string                       `bson:"name"`
	TimezoneOffset int                          `bson:"timezoneOffset"`
	Type           string                       `bson:"type"`
	Activities     []EventScheduleActivityModel `bson:"activities"`
	Date           time.Time                    `bson:"date"`
	Days           []time.Weekday               `bson:"days"`
	StartDate      time.Time                    `bson:"startDate"`
	EndDate        time.Time                    `bson:"endDate"`
}

func toEventScheduleModel(e *entities.EventSchedule) EventScheduleModel {
	return EventScheduleModel{
		ID:             e.ID,
		Name:           e.Name,
		TimezoneOffset: e.TimezoneOffset,
		Type:           string(e.Type),
		Activities: lo.Map(e.Activities, func(e entities.EventScheduleActivity, _ int) EventScheduleActivityModel {
			return toEventScheduleActivityModel(&e)
		}),
		Date:      e.Date,
		Days:      e.Days,
		StartDate: e.StartDate,
		EndDate:   e.EndDate,
	}
}

func (e *EventScheduleModel) ToEventSchedule() *entities.EventSchedule {
	return &entities.EventSchedule{
		ID:             e.ID,
		Name:           e.Name,
		TimezoneOffset: e.TimezoneOffset,
		Type:           entities.EventScheduleType(e.Type),
		Activities: lo.Map(e.Activities, func(e EventScheduleActivityModel, _ int) entities.EventScheduleActivity {
			return e.ToEventScheduleActivity()
		}),
		OneTimeEventSchedule: entities.OneTimeEventSchedule{
			Date: e.Date,
		},
		WeeklyEventSchedule: entities.WeeklyEventSchedule{
			Days: e.Days,
		},
		DailyEventSchedule: entities.DailyEventSchedule{
			StartDate: e.StartDate,
			EndDate:   e.EndDate,
		},
	}
}

type EventScheduleActivityModel struct {
	ID     string `bson:"_id"`
	Name   string `bson:"name"`
	Hour   int    `bson:"hour"`
	Minute int    `bson:"minute"`
}

func toEventScheduleActivityModel(e *entities.EventScheduleActivity) EventScheduleActivityModel {
	return EventScheduleActivityModel{
		ID:     e.ID,
		Name:   e.Name,
		Hour:   e.Hour,
		Minute: e.Minute,
	}
}

func (e *EventScheduleActivityModel) ToEventScheduleActivity() entities.EventScheduleActivity {
	return entities.EventScheduleActivity{
		ID:     e.ID,
		Name:   e.Name,
		Hour:   e.Hour,
		Minute: e.Minute,
	}
}

type EventModel struct {
	ID              string               `bson:"_id"`
	Name            string               `bson:"name"`
	ScheduleID      string               `bson:"scheduleId"`
	EventActivities []EventActivityModel `bson:"eventActivities"`
	Date            time.Time            `bson:"date"`
}

func (e *EventModel) ToEvent() *entities.Event {
	return &entities.Event{
		ID:         e.ID,
		ScheduleID: e.ScheduleID,
		Name:       e.Name,
		EventActivities: lo.Map(e.EventActivities, func(e EventActivityModel, _ int) *entities.EventActivity {
			return e.ToEventActivity()
		}),
		Date: e.Date,
	}
}

func toEventModel(e *entities.Event) EventModel {
	return EventModel{
		ID:         e.ID,
		Name:       e.Name,
		ScheduleID: e.ScheduleID,
		EventActivities: lo.Map(e.EventActivities, func(e *entities.EventActivity, _ int) EventActivityModel {
			return toEventActivityModel(e)
		}),
		Date: e.Date,
	}
}

type EventActivityModel struct {
	ID   string    `bson:"_id"`
	Name string    `bson:"name"`
	Time time.Time `bson:"time"`
}

func (e *EventActivityModel) ToEventActivity() *entities.EventActivity {
	return &entities.EventActivity{
		ID:   e.ID,
		Name: e.Name,
		Time: e.Time,
	}
}

func toEventActivityModel(e *entities.EventActivity) EventActivityModel {
	return EventActivityModel{
		ID:   e.ID,
		Name: e.Name,
		Time: e.Time,
	}
}

type PersonModel struct {
	ID                string `bson:"_id"`
	PersonID          string `bson:"personID"`
	FirstName         string `bson:"firstName"`
	MiddleName        string `bson:"middleName"`
	LastName          string `bson:"lastName"`
	ProfilePictureUrl string `bson:"profilePictureUrl"`
	Birthday          time.Time
}

func (p *PersonModel) ToEntity() *entities.Person {
	id := p.ID
	if id == "" {
		id = p.PersonID
	}
	return &entities.Person{
		PersonID:          id,
		FirstName:         p.FirstName,
		MiddleName:        p.MiddleName,
		LastName:          p.LastName,
		ProfilePictureUrl: p.ProfilePictureUrl,
	}
}

func (p *PersonModel) ToDTO() dto.PersonDTO {
	id := p.ID
	if id == "" {
		id = p.PersonID
	}
	return dto.PersonDTO{
		ID:                id,
		FirstName:         p.FirstName,
		MiddleName:        p.MiddleName,
		LastName:          p.LastName,
		ProfilePictureUrl: p.ProfilePictureUrl,
		Age:               int(time.Now().Sub(p.Birthday).Hours() / 24 / 365),
	}
}

type PersonHouseholdModel struct {
	ID          string `bson:"_id"`
	HouseholdID string `bson:"householdID"`
}

type HouseholdModel struct {
	ID               string        `bson:"_id"`
	Name             string        `bson:"name"`
	HouseholdHead    PersonModel   `bson:"householdHead"`
	PictureUrl       string        `bson:"pictureUrl"`
	HouseholdMembers []PersonModel `bson:"householdMembers"`
}

type AttendanceModel struct {
	ID            string             `bson:"_id"`
	Event         EventModel         `bson:"event"`
	EventActivity EventActivityModel `bson:"eventActivity"`

	Attendee    PersonModel `bson:"attendee"`
	CheckedInBy PersonModel `bson:"checkedInBy"`

	SecurityCode   string    `bson:"securityCode"`
	SecurityNumber int       `bson:"securityNumber"`
	CheckinTime    time.Time `bson:"checkinTime"`

	Type string `bson:"type"`
}

func toAttendanceModel(e *entities.Attendance) AttendanceModel {
	return AttendanceModel{
		ID:            e.ID,
		Event:         toEventModel(e.Event),
		EventActivity: toEventActivityModel(e.EventActivity),
		Attendee: PersonModel{
			ID:                e.Attendee.PersonID,
			PersonID:          e.Attendee.PersonID,
			FirstName:         e.Attendee.FirstName,
			MiddleName:        e.Attendee.MiddleName,
			LastName:          e.Attendee.LastName,
			ProfilePictureUrl: e.Attendee.ProfilePictureUrl,
			Birthday:          time.Time{},
		},
		CheckedInBy: PersonModel{
			ID:                e.CheckedInBy.PersonID,
			PersonID:          e.CheckedInBy.PersonID,
			FirstName:         e.CheckedInBy.FirstName,
			MiddleName:        e.CheckedInBy.MiddleName,
			LastName:          e.CheckedInBy.LastName,
			ProfilePictureUrl: e.CheckedInBy.ProfilePictureUrl,
			Birthday:          time.Time{},
		},
		SecurityCode:   e.SecurityCode,
		SecurityNumber: e.SecurityNumber,
		CheckinTime:    e.CheckinTime,
		Type:           string(e.Type),
	}
}

func (e *AttendanceModel) ToAttendance() *entities.Attendance {
	return &entities.Attendance{
		ID:             e.ID,
		Event:          e.Event.ToEvent(),
		EventActivity:  e.EventActivity.ToEventActivity(),
		Attendee:       &entities.Person{PersonID: e.Attendee.PersonID, FirstName: e.Attendee.FirstName, MiddleName: e.Attendee.MiddleName, LastName: e.Attendee.LastName, ProfilePictureUrl: e.Attendee.ProfilePictureUrl},
		SecurityCode:   e.SecurityCode,
		SecurityNumber: e.SecurityNumber,
		CheckinTime:    e.CheckinTime,
		Type:           entities.AttendanceType(e.Type),
	}
}

type ActivityAttendanceSummaryModel struct {
	ID          string         `bson:"_id"`
	Name        string         `bson:"Name"`
	Total       int            `bson:"total"`
	TotalByType map[string]int `bson:"totalByType"`
}

func (e *ActivityAttendanceSummaryModel) ToDTO() dto.ActivityAttendanceSummaryDTO {
	return dto.ActivityAttendanceSummaryDTO{
		Name:        e.Name,
		Total:       e.Total,
		TotalByType: e.TotalByType,
	}
}

type EventAttendanceSummaryModel struct {
	ID              string    `bson:"_id"`
	Date            time.Time `bson:"date"`
	TotalCheckedIn  int       `bson:"totalCheckedIn"`
	TotalCheckedOut int       `bson:"totalCheckedOut"`
	TotalFirstTimer int       `bson:"totalFirstTimer"`
	Total           int       `bson:"total"`

	TotalByType        map[string]int                   `bson:"totalByType"`
	AcitivitiesSummary []ActivityAttendanceSummaryModel `bson:"activitiesSummary"`
	LastUpdated        time.Time                        `bson:"lastUpdated"`
	NextUpdate         time.Time                        `bson:"nextUpdate"`
}

func (e *EventAttendanceSummaryModel) ToDTO() dto.EventAttendanceSummaryDTO {
	return dto.EventAttendanceSummaryDTO{
		TotalCheckedIn:  e.TotalCheckedIn,
		TotalCheckedOut: e.TotalCheckedOut,
		TotalFirstTimer: e.TotalFirstTimer,
		Total:           e.Total,
		TotalByType:     e.TotalByType,
		AcitivitiesSummary: lo.Map(e.AcitivitiesSummary, func(ee ActivityAttendanceSummaryModel, _ int) dto.ActivityAttendanceSummaryDTO {
			return ee.ToDTO()
		}),
		Date: e.Date,
		ID:   e.ID,
	}
}
