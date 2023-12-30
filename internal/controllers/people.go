package controllers

import (
	"github.com/Goldwin/ies-pik-cms/internal/middleware"
	"github.com/Goldwin/ies-pik-cms/pkg/common/commands"
	"github.com/Goldwin/ies-pik-cms/pkg/people"
	"github.com/Goldwin/ies-pik-cms/pkg/people/dto"
	"github.com/gin-gonic/gin"
)

type peopleManagementController struct {
	peopleComponent people.PeopleManagementComponent
	middleware      middleware.MiddlewareComponent
}

func InitializePeopleManagementController(
	r *gin.Engine,
	middlewareComponent middleware.MiddlewareComponent,
	peopleComponent people.PeopleManagementComponent,
) {
	c := &peopleManagementController{
		peopleComponent: peopleComponent,
		middleware:      middlewareComponent,
	}
	rg := r.Group("people")
	rg.POST("person", middlewareComponent.Auth("PERSON_ADD"), c.addPersonInfo)
	rg.PUT("person/:id", middlewareComponent.Auth("PERSON_UPDATE"), c.updatePersonInfo)
	rg.POST("household", middlewareComponent.Auth("HOUSEHOLD_ADD"), c.addHousehold)
	rg.PUT("household/:id", middlewareComponent.Auth("HOUSEHOLD_UPDATE"), c.updateHousehold)
}

func (c *peopleManagementController) addPersonInfo(ctx *gin.Context) {
	var input dto.Person
	ctx.BindJSON(&input)
	input.ID = ""
	c.peopleComponent.AddPerson(ctx, input, &outputDecorator[dto.Person]{
		output: nil,
		errFunction: func(err commands.AppErrorDetail) {
			ctx.AbortWithStatusJSON(400, gin.H{
				"error": err.Error(),
			})
		},
		successFunc: func(person dto.Person) {
			ctx.JSON(200, gin.H{
				"data": person,
			})
		},
	})
}

func (c *peopleManagementController) addHousehold(ctx *gin.Context) {
	var input dto.HouseHoldInput
	ctx.BindJSON(&input)
	input.ID = ""
	c.peopleComponent.AddHousehold(ctx, input, &outputDecorator[dto.Household]{
		output: nil,
		errFunction: func(err commands.AppErrorDetail) {
			ctx.AbortWithStatusJSON(400, gin.H{
				"error": err.Error(),
			})
		},
		successFunc: func(person dto.Household) {
			ctx.JSON(200, gin.H{
				"data": person,
			})
		},
	})
}

func (c *peopleManagementController) updateHousehold(ctx *gin.Context) {
	var input dto.HouseHoldInput
	ctx.BindJSON(&input)
	input.ID = ctx.Param("id")
	c.peopleComponent.UpdateHousehold(ctx, input, &outputDecorator[dto.Household]{
		output: nil,
		errFunction: func(err commands.AppErrorDetail) {
			ctx.AbortWithStatusJSON(400, gin.H{
				"error": err.Error(),
			})
		},
		successFunc: func(person dto.Household) {
			ctx.JSON(200, gin.H{
				"data": person,
			})
		},
	})
}

func (c *peopleManagementController) updatePersonInfo(ctx *gin.Context) {
	var input dto.Person
	ctx.BindJSON(&input)
	input.ID = ctx.Param("id")
	c.peopleComponent.UpdatePerson(ctx, input, &outputDecorator[dto.Person]{
		output: nil,
		errFunction: func(err commands.AppErrorDetail) {
			ctx.AbortWithStatusJSON(400, gin.H{
				"error": err.Error(),
			})
		},
		successFunc: func(person dto.Person) {
			ctx.JSON(200, gin.H{
				"data": person,
			})
		},
	})
}
