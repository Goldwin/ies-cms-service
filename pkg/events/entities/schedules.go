package entities

import (
	"time"
)

type ChurchEventSchedule struct {
	ID             string
	DayOfWeek      int
	Hours          int
	Minute         int
	TimezoneOffset int
}

func (e ChurchEventSchedule) getNextWeekEventDate() time.Time {
	days := e.DayOfWeek - int(time.Now().UTC().Add(time.Hour*time.Duration(e.TimezoneOffset)).Weekday())
	y, m, d := time.Now().AddDate(0, 0, days).Date()
	return time.Date(y, m, d, e.Hours, e.Minute, 0, 0, time.UTC).Add(time.Hour * time.Duration(-e.TimezoneOffset))
}
