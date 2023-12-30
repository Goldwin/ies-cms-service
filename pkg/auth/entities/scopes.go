package entities

type Scope string

type Role struct {
	ID     int
	Name   string
	Scopes []Scope
}

const (
	AllAccess Scope = "ALL_ACCESS"

	EventView    Scope = "EVENT_VIEW"
	EventCheckIn Scope = "EVENT_CHECK_IN"
)

var (
	ChurchMember = Role{
		ID:   1,
		Name: "Church Member",
		Scopes: []Scope{
			EventCheckIn,
			EventView,
		},
	}
	Admin = Role{
		ID:   0,
		Name: "Admin",
		Scopes: []Scope{
			AllAccess,
		},
	}
)
