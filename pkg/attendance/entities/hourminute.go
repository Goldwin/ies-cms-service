package entities

import "fmt"

type HourMinute struct {
	Hour   int
	Minute int
}

func (h *HourMinute) IsValid() string {
	if h.Hour < 0 || h.Hour > 23 {
		return "Hour must be between 0 and 23"
	}
	if h.Minute < 0 || h.Minute > 59 {
		return "Minute must be between 0 and 59"
	}
	return ""
}

func (h *HourMinute) String() string {
	return fmt.Sprintf("%02d:%02d", h.Hour, h.Minute)
}

func (h *HourMinute) SetFromString(s string) error {
	_, err := fmt.Sscanf(s, "%d:%d", &h.Hour, &h.Minute)
	if err != nil {
		return fmt.Errorf("invalid time: %s", s)
	}
	return nil
}

func (h *HourMinute) SetFromStringOrZero(s string) {
	if h.SetFromString(s) != nil {
		h.Hour = 0
		h.Minute = 0
	}
}

func (h *HourMinute) SetFromStringOrMaxValue(s string) { 
	if h.SetFromString(s) != nil {
		h.Hour = 23
		h.Minute = 59
	}
}