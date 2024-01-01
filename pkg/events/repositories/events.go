package repositories

import "github.com/Goldwin/ies-pik-cms/pkg/events/entities"

type ChurchEventRepository interface {
	Save(entities.ChurchEvent) error
	Get(string) (*entities.ChurchEvent, error)
}
