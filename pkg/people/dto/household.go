package dto

type HouseholdPerson struct {
	ID                string `json:"id"`
	FirstName         string `json:"first_name"`
	MiddleName        string `json:"middle_name"`
	LastName          string `json:"last_name"`
	PhoneNumber       string `json:"phone_number"`
	EmailAddress      string `json:"email_address"`
	ProfilePictureUrl string `json:"profile_picture_url"`
}

type Household struct {
	ID            string            `json:"id"`
	Name          string            `json:"name"`
	PictureUrl    string            `json:"picture_url"`
	HouseholdHead HouseholdPerson   `json:"household_head"`
	Members       []HouseholdPerson `json:"members"`
}

type HouseHoldInput struct {
	ID                    string   `json:"id"`
	Name                  string   `json:"name"`
	HouseholdHeadPersonId string   `json:"head_person_id"`
	MemberPersonsIds      []string `json:"member_person_ids"`
}
