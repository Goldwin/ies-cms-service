package queries

type QueryContext interface {
	GetEventSchedule()
	ListEventSchedules() ListEventSchedule
	ListEventBySchedule()
	GetEvent()
	ListEventCheckIn()
}
