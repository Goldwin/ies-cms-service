package repositories

import "github.com/Goldwin/ies-pik-cms/pkg/people/queries"

type QueryContext interface {
	SearchPerson() queries.SearchPerson
	ViewPerson() queries.ViewPerson
}
