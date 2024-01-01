package repositories

import (
	"time"

	"github.com/Goldwin/ies-pik-cms/pkg/events/entities"
)

type ChurchEventScheduleRepository interface {
	Save(entities.ChurchEventSchedule) error
	GetByTimezoneAndWeekDay(timezoneOffset int, weekday time.Weekday) (*entities.ChurchEventSchedule, error)
}
