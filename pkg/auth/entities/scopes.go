package entities

type Scope string

type Role struct {
	ID     int
	Name   string
	Scopes []Scope
}

const (
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
)
