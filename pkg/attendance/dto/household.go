package dto

type PersonDTO struct {
	ID                string `json:"id"`
	FirstName         string `json:"firstName"`
	MiddleName        string `json:"middleName"`
	LastName          string `json:"lastName"`
	ProfilePictureUrl string `json:"profilePictureUrl"`
	Age               int    `json:"age"`
}

type HouseholdDTO struct {
	ID            string      `json:"id"`
	Name          string      `json:"name"`
	PictureUrl    string      `json:"pictureUrl"`
	HouseholdHead PersonDTO   `json:"householdHead"`
	Members       []PersonDTO `json:"members"`
}
