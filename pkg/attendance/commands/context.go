package commands

import "github.com/Goldwin/ies-pik-cms/pkg/attendance/repositories"

type CommandContext interface {
	EventRepository() repositories.EventRepository
	EventScheduleRepository() repositories.EventScheduleRepository
	AttendanceRepository() repositories.AttendanceRepository
}