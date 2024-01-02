package controllers

import (
	"github.com/Goldwin/ies-pik-cms/internal/bus"
	"github.com/Goldwin/ies-pik-cms/internal/middleware"
	"github.com/Goldwin/ies-pik-cms/pkg/common/commands"
	"github.com/Goldwin/ies-pik-cms/pkg/events"
	"github.com/Goldwin/ies-pik-cms/pkg/events/dto"
	"github.com/gin-gonic/gin"
)

type churchEventController struct {
	churchEventComponent events.ChurchEventComponent
	middleware           middleware.MiddlewareComponent
}

func InitializeEventsController(
	r *gin.Engine,
	middlewareComponent middleware.MiddlewareComponent,
	peopleComponent events.ChurchEventComponent,
	eventBusComponent bus.EventBusComponent,
) {
	c := &churchEventController{
		churchEventComponent: peopleComponent,
		middleware:           middlewareComponent,
	}

	rg := r.Group("events")
	rg.POST("checkin", middlewareComponent.Auth("EVENT_CHECK_IN"), c.checkIn)
}

func (c *churchEventController) checkIn(ctx *gin.Context) {
	input := dto.CheckInInput{}
	ctx.BindJSON(&input)
	c.churchEventComponent.CheckIn(ctx, input, &outputDecorator[[]dto.CheckInEvent]{
		output: nil,
		errFunction: func(err commands.AppErrorDetail) {
			ctx.AbortWithStatusJSON(400, gin.H{
				"error": err.Error(),
			})
		},
		successFunc: func(response []dto.CheckInEvent) {
			ctx.JSON(200, gin.H{
				"data": response,
			})
		},
	})
}
