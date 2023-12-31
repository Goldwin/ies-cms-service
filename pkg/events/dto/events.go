package dto

type CreateEventInput struct {
	ID     string
	Year   int
	Month  int
	Day    int
	Hours  int
	Minute int
}

type ChurchEvent struct {
	ID     string
	Year   int
	Month  int
	Day    int
	Hours  int
	Minute int
}
