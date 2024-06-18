package controllers

import (
	"github.com/Goldwin/ies-pik-cms/internal/middleware"
	"github.com/Goldwin/ies-pik-cms/pkg/attendance"
	"github.com/Goldwin/ies-pik-cms/pkg/common/out"
	"github.com/gin-gonic/gin"
)

type attendanceController struct {
	attendanceComponent attendance.AttendanceComponent
	middlewareComponent middleware.MiddlewareComponent
}

func InitializeAttendanceController(r *gin.Engine, middleware middleware.MiddlewareComponent, attendanceComponent attendance.AttendanceComponent) {
	attendanceController := attendanceController{
		attendanceComponent: attendanceComponent,
		middlewareComponent: middleware,
	}
	
	rg := r.Group("attendance")
	rg.POST("schedules", attendanceController.createSchedule)
	
	scheduleURL := "schedules/:id"
	rg.GET(scheduleURL)
	rg.PUT(scheduleURL)
	rg.DELETE(scheduleURL)

	rg.GET("schedules/:scheduleID/events")
	rg.GET("schedules/:scheduleID/events/:eventID")

	rg.GET("schedules/:scheduleID/events/:eventID/checkin")
	rg.POST("schedules/:scheduleID/events/:eventID/checkin")

	rg.GET("schedules/:scheduleID/stats")
}

func (a *attendanceController) createSchedule(c *gin.Context) {
	var schedule attendance.EventScheduleDTO
	output := &outputDecorator[attendance.EventScheduleDTO]{
		output: nil,
		errFunction: func(err out.AppErrorDetail) {
			c.JSON(400, gin.H{
				"error": err,
			})
		},
		successFunc: func(result attendance.EventScheduleDTO) {
			c.JSON(200, gin.H{
				"data": result,
			})
		},
	}
	a.attendanceComponent.CreateEventSchedule(c, schedule, output)
}
