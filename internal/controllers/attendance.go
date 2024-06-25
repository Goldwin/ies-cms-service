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
	rg.GET(scheduleURL, attendanceController.getEventSchedule)
	rg.PUT(scheduleURL, attendanceController.updateEventSchedule)
	rg.DELETE(scheduleURL, attendanceController.deleteEventSchedule)

	rg.GET("schedules/:scheduleID/events", attendanceController.listEventsBySchedule)
	rg.GET("schedules/:scheduleID/events/:eventID", attendanceController.getEventBySchedule)

	rg.GET("schedules/:scheduleID/events/:eventID/checkin", attendanceController.listEventCheckIn)
	rg.POST("schedules/:scheduleID/events/:eventID/checkin", attendanceController.checkIn)

	rg.GET("schedules/:scheduleID/stats", attendanceController.getEventScheduleStats)
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

func (a *attendanceController) getEventSchedule(c *gin.Context) {
	//TODO fill this
	// id := c.Param("id")
	// output := &outputDecorator[attendance.EventScheduleDTO]{
	// 	output: nil,
	// 	errFunction: func(err out.AppErrorDetail) {
	// 		c.JSON(400, gin.H{
	// 			"error": err,
	// 		})
	// 	},
	// 	successFunc: func(result attendance.EventScheduleDTO) {
	// 		c.JSON(200, gin.H{
	// 			"data": result,
	// 		})
	// 	},
	// }
	// //a.attendanceComponent.GetEventSchedule(c, id, output)
}

func (a *attendanceController) updateEventSchedule(c *gin.Context) {
	//TODO fill this
}

func (a *attendanceController) deleteEventSchedule(c *gin.Context) {
	//TODO fill this
}	

func (a *attendanceController) getEventScheduleStats(c *gin.Context) {
	//TODO fill this
}

func (a *attendanceController) listEventsBySchedule(c *gin.Context) {
	//TODO fill this
}	

func (a *attendanceController) getEventBySchedule(c *gin.Context) {
	//TODO fill this
}

func (a *attendanceController) listEventCheckIn(c *gin.Context) {
	//TODO fill this
}

func (a *attendanceController) checkIn(c *gin.Context) {
	//TODO fill this
}


