package dto

type ChurchEventSchedule struct {
	ID             string
	Name           string
	DayOfWeek      int
	Hours          int
	Minute         int
	TimezoneOffset int
}
