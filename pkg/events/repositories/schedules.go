package repositories

import "github.com/Goldwin/ies-pik-cms/pkg/events/dto"

type ChurchEventScheduleRepository interface {
	Save(dto.ChurchEventSchedule) error
	Get(string) (*dto.ChurchEventSchedule, error)
}
