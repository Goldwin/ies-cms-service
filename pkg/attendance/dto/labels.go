package dto

import (
	"github.com/Goldwin/ies-pik-cms/pkg/attendance/entities"
	"github.com/samber/lo"
)

type LabelDTO struct {
	ID          string           `json:"id"`
	Name        string           `json:"name"`
	Type        string           `json:"type"`
	Orientation string           `json:"orientation"`
	PaperSize   []float64        `json:"paperSize,omitempty" `
	Margin      []float64        `json:"margin,omitempty"`
	Objects     []map[string]any `json:"objects,omitempty"`
}

type ActivityLabelDTO struct {
	LabelID         string   `json:"labelId"`
	LabelName       string   `json:"labelName"`
	Type            string   `json:"type"`
	AttendanceTypes []string `json:"attendanceTypes"`
	Quantity        int      `json:"quantity"`
}

func (d ActivityLabelDTO) ToEntity() *entities.ActivityLabel {
	return &entities.ActivityLabel{
		LabelID:   d.LabelID,
		LabelName: d.LabelName,
		AttendanceTypes: lo.Map(d.AttendanceTypes, func(attendanceType string, _ int) entities.AttendanceType {
			return entities.AttendanceType(attendanceType)
		}),
		Quantity: d.Quantity,
	}
}

func FromActivityLabelEntity(label *entities.ActivityLabel) ActivityLabelDTO {
	return ActivityLabelDTO{
		LabelID:   label.LabelID,
		LabelName: label.LabelName,
		AttendanceTypes: lo.Map(label.AttendanceTypes, func(attendanceType entities.AttendanceType, _ int) string {
			return string(attendanceType)
		}),
		Quantity: label.Quantity,
	}
}
