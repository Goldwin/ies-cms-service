package entities

type Scope string

type Role struct {
	ID     string
	Name   string
	Scopes []Scope
}

const (
	AllAccess Scope = "ALL_ACCESS"

	EventView    Scope = "EVENT_VIEW"
	EventCheckIn Scope = "EVENT_CHECK_IN"

	ProfileView   Scope = "PROFILE_VIEW"
	ProfileUpdate Scope = "PROFILE_UPDATE"
)

var (
	ChurchMember = Role{
		ID:   "1",
		Name: "Church Member",
		Scopes: []Scope{
			EventCheckIn,
			EventView,
			ProfileView,
			ProfileUpdate,
		},
	}
	Admin = Role{
		ID:   "2",
		Name: "Admin",
		Scopes: []Scope{
			AllAccess,
		},
	}
)
