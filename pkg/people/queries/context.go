package queries

type QueryContext interface {
	SearchPerson() SearchPerson
	ViewPerson() ViewPerson
	ViewPersonByEmail() ViewPersonByEmail
	ViewHouseholdByPerson() ViewHouseholdByPerson
}
