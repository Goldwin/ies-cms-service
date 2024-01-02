package repositories

import "github.com/Goldwin/ies-pik-cms/pkg/events/entities"

type ChurchEventSessionRepository interface {
	Save(entities.ChurchEventSession) error
	Get(ID string) (*entities.ChurchEventSession, error)
}
