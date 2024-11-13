package entities

type Label struct {
	ID          string
	Name        string
	Orientation string
	PaperSize   []float64
	Margin      []float64
	Objects     []map[string]any
}

type ActivityLabel struct {
	LabelID         string
	LabelName       string
	AttendanceTypes []AttendanceType
	Quantity        int
}
