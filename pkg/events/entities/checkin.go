package entities

type CheckInEvent struct {
	ID     string
	Person Person
	Event  ChurchEvent
}
