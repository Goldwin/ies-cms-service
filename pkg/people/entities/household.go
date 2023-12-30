package entities

type Household struct {
	ID            string
	Name          string
	HouseholdHead Person
	PictureUrl    string
	Members       []Person
}
