package entities

type LabelType string

const (
	NameLabelType     LabelType = "name"
	SecurityLabelType LabelType = "security"
)

type Label struct {
	ID          string
	Name        string
	Type        LabelType
	Orientation string
	PaperSize   []float64
	Margin      []float64
	Objects     []map[string]any
}

type ActivityLabel struct {
	LabelID         string
	LabelName       string
	Type            LabelType
	AttendanceTypes []AttendanceType
	Quantity        int
}
