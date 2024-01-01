package repositories

import "github.com/Goldwin/ies-pik-cms/pkg/events/entities"

type ChurchEventScheduleRepository interface {
	Save(entities.ChurchEventSchedule) error
	Get(string) (*entities.ChurchEventSchedule, error)
}
