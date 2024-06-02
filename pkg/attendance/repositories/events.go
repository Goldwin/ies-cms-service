package repositories

import (
	"github.com/Goldwin/ies-pik-cms/pkg/attendance/entities"
	. "github.com/Goldwin/ies-pik-cms/pkg/common/repositories"
)

type EventRepository interface {
	Repository[string, entities.Event]
}
