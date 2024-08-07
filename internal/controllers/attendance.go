package controllers

import (
	"github.com/Goldwin/ies-pik-cms/internal/middleware"
	"github.com/Goldwin/ies-pik-cms/pkg/attendance"
	"github.com/Goldwin/ies-pik-cms/pkg/attendance/dto"
	"github.com/Goldwin/ies-pik-cms/pkg/attendance/queries"
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
	rg.GET("schedules", attendanceController.listEventSchedules)

	scheduleURL := "schedules/:scheduleID"
	rg.GET(scheduleURL, attendanceController.getEventSchedule)
	rg.PUT(scheduleURL, attendanceController.updateEventSchedule)
	rg.DELETE(scheduleURL, attendanceController.archiveEventSchedule)

	activitiesUrl := scheduleURL + "/activities"
	activityUrl := activitiesUrl + "/:activityID"
	rg.POST(activitiesUrl, attendanceController.createEventScheduleActivity)
	rg.PUT(activityUrl, attendanceController.updateEventScheduleActivity)
	rg.DELETE(activityUrl, attendanceController.removeEventScheduleActivity)

	rg.GET("schedules/:scheduleID/events", attendanceController.listEventsBySchedule)
	rg.GET("schedules/:scheduleID/events/:eventID", attendanceController.getEventBySchedule)

	rg.GET("schedules/:scheduleID/events/:eventID/checkin", attendanceController.listEventCheckIn)
	rg.POST("schedules/:scheduleID/events/:eventID/checkin", attendanceController.checkIn)

	rg.GET("schedules/:scheduleID/stats", attendanceController.getEventScheduleStats)
}

func (a *attendanceController) listEventSchedules(c *gin.Context) {
	var query queries.ListEventScheduleQuery

	err := c.ShouldBindQuery(&query)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	output := &outputDecorator[queries.ListEventScheduleResult]{
		output: nil,
		errFunction: func(err out.AppErrorDetail) {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
		},
		successFunc: func(result queries.ListEventScheduleResult) {
			c.JSON(200, result)
		},
	}
	a.attendanceComponent.ListEventSchedules(c, query, output)
}

func (a *attendanceController) createSchedule(c *gin.Context) {
	var schedule dto.EventScheduleDTO
	output := &outputDecorator[dto.EventScheduleDTO]{
		output: nil,
		errFunction: func(err out.AppErrorDetail) {
			c.JSON(400, gin.H{
				"error": err,
			})
		},
		successFunc: func(result dto.EventScheduleDTO) {
			c.JSON(200, gin.H{
				"data": result,
			})
		},
	}
	err := c.ShouldBindJSON(&schedule)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	a.attendanceComponent.CreateEventSchedule(c, schedule, output).Wait()
}

func (a *attendanceController) getEventSchedule(c *gin.Context) {
	id := c.Param("scheduleID")
	output := &outputDecorator[queries.GetEventScheduleResult]{
		output: nil,
		errFunction: func(err out.AppErrorDetail) {
			c.JSON(400, gin.H{
				"error": err,
			})
		},
		successFunc: func(result queries.GetEventScheduleResult) {
			c.JSON(200, result)
		},
	}
	a.attendanceComponent.GetEventSchedule(c, queries.GetEventScheduleQuery{
		ScheduleID: id,
	}, output)
}

func (a *attendanceController) updateEventSchedule(c *gin.Context) {
	var schedule dto.EventScheduleDTO
	output := &outputDecorator[dto.EventScheduleDTO]{
		output: nil,
		errFunction: func(err out.AppErrorDetail) {
			c.JSON(400, gin.H{
				"error": err,
			})
		},
		successFunc: func(result dto.EventScheduleDTO) {
			c.JSON(200, gin.H{
				"data": result,
			})
		},
	}
	err := c.ShouldBindJSON(&schedule)
	schedule.ID = c.Param("scheduleID")
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	a.attendanceComponent.UpdateEventSchedule(c, schedule, output).Wait()
}

func (a *attendanceController) archiveEventSchedule(c *gin.Context) {
	//TODO fill this
}

func (a *attendanceController) getEventScheduleStats(c *gin.Context) {
	//TODO fill this
}

func (a *attendanceController) listEventsBySchedule(c *gin.Context) {
	var input queries.ListEventByScheduleQuery
	err := c.BindQuery(&input)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	output := &outputDecorator[queries.ListEventByScheduleResult]{
		output: nil,
		errFunction: func(err out.AppErrorDetail) {
			c.JSON(400, err)
		},
		successFunc: func(result queries.ListEventByScheduleResult) {
			c.JSON(200, result)
		},
	}
	a.attendanceComponent.ListEventsBySchedule(c, input, output)
}

func (a *attendanceController) getEventBySchedule(c *gin.Context) {
	var input queries.GetEventQuery

	err := c.BindQuery(&input)

	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	output := &outputDecorator[queries.GetEventResult]{
		output: nil,
		errFunction: func(err out.AppErrorDetail) {
			c.JSON(400, err)
		},
		successFunc: func(result queries.GetEventResult) {
			c.JSON(200, result)
		},
	}

	a.attendanceComponent.GetEvent(c, input, output)
}

func (a *attendanceController) listEventCheckIn(c *gin.Context) {
	var input queries.ListEventAttendanceQuery
	err := c.ShouldBindQuery(&input)

	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	output := &outputDecorator[queries.ListEventAttendanceResult]{
		output: nil,
		errFunction: func(err out.AppErrorDetail) {
			c.JSON(400, err)
		},
		successFunc: func(result queries.ListEventAttendanceResult) {
			c.JSON(200, result)
		},
	}

	a.attendanceComponent.ListEventAttendance(c, input, output)
}

func (a *attendanceController) createEventScheduleActivity(c *gin.Context) {
	var data dto.EventScheduleActivityDTO
	err := c.ShouldBindJSON(&data)
	data.ScheduleID = c.Param("scheduleID")

	if err != nil {
		c.JSON(400, gin.H{
			"error": gin.H{
				"code":    "400",
				"message": err.Error(),
			},
		})
		return
	}

	output := &outputDecorator[dto.EventScheduleDTO]{
		output: nil,
		errFunction: func(err out.AppErrorDetail) {
			c.JSON(400, gin.H{
				"error": err,
			})
		},
		successFunc: func(result dto.EventScheduleDTO) {
			c.JSON(200, gin.H{
				"data": result,
			})
		},
	}

	a.attendanceComponent.AddEventScheduleActivity(c, data, output).Wait()
}

func (a *attendanceController) updateEventScheduleActivity(c *gin.Context) {
	//TODO fill this
}

func (a *attendanceController) removeEventScheduleActivity(c *gin.Context) {
	//TODO fill this
}

func (a *attendanceController) checkIn(c *gin.Context) {
	//TODO fill this
}
