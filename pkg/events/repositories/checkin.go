package repositories

import "github.com/Goldwin/ies-pik-cms/pkg/events/entities"

type EventCheckInRepository interface {
	Save(entities.CheckInEvent) error
	Get(string) (*entities.CheckInEvent, error)
}
