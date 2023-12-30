package dto

type HouseholdPerson struct {
	ID           string
	FirstName    string
	MiddleName   string
	LastName     string
	PhoneNumber  string
	EmailAddress string
}

type Household struct {
	ID            string
	Name          string
	PictureUrl    string
	HouseholdHead HouseholdPerson
	Members       []HouseholdPerson
}

type HouseHoldInput struct {
	ID                    string
	Name                  string
	HouseholdHeadPersonId string
	PictureUrl            string
	MemberPersonsIds      []string
}
