package queries

type QueryContext interface {
	GetEventSchedule() GetEventSchedule
	ListEventSchedules() ListEventSchedule
	ListEventsBySchedule() ListEventBySchedule
	GetEvent() GetEvent
	ListEventAttendance() ListEventAttendance
	GetEventAttendanceSummary() GetEventAttendanceSummary
	GetEventScheduleStats() GetEventScheduleStats
	ListLabels() ListLabels
}
