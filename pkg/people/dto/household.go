package dto

type HouseholdPerson struct {
	ID                string `json:"id"`
	FirstName         string `json:"firstName"`
	MiddleName        string `json:"middleName"`
	LastName          string `json:"lastName"`
	PhoneNumber       string `json:"phoneNumber"`
	EmailAddress      string `json:"emailAddress"`
	ProfilePictureUrl string `json:"profilePictureUrl"`
}

type Household struct {
	ID            string            `json:"id"`
	Name          string            `json:"name"`
	PictureUrl    string            `json:"pictureUrl"`
	HouseholdHead HouseholdPerson   `json:"householdHead"`
	Members       []HouseholdPerson `json:"members"`
}

type HouseHoldInput struct {
	ID                    string   `json:"id"`
	Name                  string   `json:"name"`
	HouseholdHeadPersonId string   `json:"headPersonId"`
	MemberPersonsIds      []string `json:"memberPersonIds"`
}
