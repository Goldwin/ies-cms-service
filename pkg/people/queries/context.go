package queries

type QueryContext interface {
	SearchPerson() SearchPerson
	ViewPerson() ViewPerson
	ViewHouseholdByPerson() ViewHouseholdByPerson
}
