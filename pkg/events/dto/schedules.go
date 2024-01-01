package dto

type ChurchEventSchedule struct {
	ID             string
	DayOfWeek      int
	Hours          int
	Minute         int
	TimezoneOffset int
}
