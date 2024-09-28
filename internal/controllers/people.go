package controllers

import (
	"github.com/Goldwin/ies-pik-cms/internal/bus"
	"github.com/Goldwin/ies-pik-cms/internal/middleware"
	"github.com/Goldwin/ies-pik-cms/pkg/common/out"
	"github.com/Goldwin/ies-pik-cms/pkg/people"
	"github.com/Goldwin/ies-pik-cms/pkg/people/dto"
	"github.com/Goldwin/ies-pik-cms/pkg/people/queries"
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
	eventBusComponent bus.EventBusComponent,
) {
	c := &peopleManagementController{
		peopleComponent: peopleComponent,
		middleware:      middlewareComponent,
	}
	personIdUrl := "persons/:id"
	rg := r.Group("people")
	rg.POST("persons", middlewareComponent.Auth("PERSON_CREATE"), c.addPersonInfo)
	rg.PUT(personIdUrl, middlewareComponent.Auth("PERSON_UPDATE"), c.updatePersonInfo)
	rg.DELETE(personIdUrl, middlewareComponent.Auth("PERSON_DELETE"), c.deletePersonInfo)
	rg.GET(personIdUrl, middlewareComponent.Auth("PERSON_VIEW"), c.viewPerson)
	rg.GET("persons/:id/household", middlewareComponent.Auth("PERSON_VIEW"), c.viewPersonHousehold)
	rg.POST("search", middlewareComponent.Auth("PERSON_SEARCH"), c.searchPerson)
	rg.POST("households", middlewareComponent.Auth("HOUSEHOLD_CREATE"), c.addHousehold)
	rg.PUT("households/:id", middlewareComponent.Auth("HOUSEHOLD_UPDATE"), c.updateHousehold)
	rg.DELETE("households/:id", middlewareComponent.Auth("HOUSEHOLD_DELETE"), c.deleteHousehold)
	rg.POST("households/search", middlewareComponent.Auth("HOUSEHOLDS_VIEW"), c.householdSearch)
}

func (c *peopleManagementController) addPersonInfo(ctx *gin.Context) {
	var input dto.Person
	err := ctx.BindJSON(&input)
	if err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	input.ID = ""
	c.peopleComponent.AddPerson(ctx, input, &outputDecorator[dto.Person]{
		output: nil,
		errFunction: func(err out.AppErrorDetail) {
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
		errFunction: func(err out.AppErrorDetail) {
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
		errFunction: func(err out.AppErrorDetail) {
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

func (c *peopleManagementController) deleteHousehold(ctx *gin.Context) {
	var input dto.HouseHoldInput
	input.ID = ctx.Param("id")
	c.peopleComponent.DeleteHousehold(ctx, input, &outputDecorator[bool]{
		output: nil,
		errFunction: func(err out.AppErrorDetail) {
			ctx.AbortWithStatusJSON(400, gin.H{
				"error": err.Error(),
			})
		},
		successFunc: func(isSuccess bool) {
			ctx.JSON(204, gin.H{})
		},
	})
}

func (c *peopleManagementController) updatePersonInfo(ctx *gin.Context) {
	var input dto.Person
	ctx.BindJSON(&input)
	input.ID = ctx.Param("id")
	c.peopleComponent.UpdatePerson(ctx, input, &outputDecorator[dto.Person]{
		output: nil,
		errFunction: func(err out.AppErrorDetail) {
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

func (c *peopleManagementController) deletePersonInfo(ctx *gin.Context) {
	var input dto.Person
	input.ID = ctx.Param("id")
	c.peopleComponent.DeletePerson(ctx, input, &outputDecorator[bool]{
		output: nil,
		errFunction: func(err out.AppErrorDetail) {
			ctx.AbortWithStatusJSON(400, gin.H{
				"error": err.Error(),
			})
		},
		successFunc: func(person bool) {
			ctx.JSON(204, gin.H{})
		},
	})
}

func (c *peopleManagementController) viewPerson(ctx *gin.Context) {
	id := ctx.Param("id")
	c.peopleComponent.ViewPerson(ctx, queries.ViewPersonQuery{
		ID: id,
	}, &outputDecorator[queries.ViewPersonResult]{
		output: nil,
		errFunction: func(err out.AppErrorDetail) {
			ctx.AbortWithStatusJSON(err.Code, gin.H{
				"error": err.Error(),
			})
		},
		successFunc: func(result queries.ViewPersonResult) {
			ctx.JSON(200, result)
		},
	})
}

func (c *peopleManagementController) viewPersonHousehold(ctx *gin.Context) {
	id := ctx.Param("id")
	c.peopleComponent.ViewHouseholdByPerson(ctx, queries.ViewHouseholdByPersonQuery{
		PersonID: id,
	}, &outputDecorator[queries.ViewHouseholdByPersonResult]{
		output: nil,
		errFunction: func(err out.AppErrorDetail) {
			ctx.AbortWithStatusJSON(err.Code, gin.H{
				"error": err.Error(),
			})
		},
		successFunc: func(result queries.ViewHouseholdByPersonResult) {
			ctx.JSON(200, result)
		},
	})
}

func (c *peopleManagementController) searchPerson(ctx *gin.Context) {
	var input queries.SearchPersonQuery
	err := ctx.BindJSON(&input)
	if err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.peopleComponent.SearchPerson(ctx, input, &outputDecorator[queries.SearchPersonResult]{
		output: nil,
		errFunction: func(err out.AppErrorDetail) {
			ctx.AbortWithStatusJSON(err.Code, gin.H{
				"error": err.Error(),
			})
		},
		successFunc: func(result queries.SearchPersonResult) {
			ctx.JSON(200, result)
		},
	})
}

func (a *peopleManagementController) householdSearch(c *gin.Context) {
	var data queries.SearchHouseholdFilter
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.JSON(400, gin.H{
			"error": gin.H{
				"code":    "400",
				"message": err.Error(),
			},
		})
	}

	output := &outputDecorator[queries.SearchHouseholdResult]{
		output: nil,
		errFunction: func(err out.AppErrorDetail) {
			c.JSON(400, gin.H{
				"error": err,
			})
		},
		successFunc: func(result queries.SearchHouseholdResult) {
			c.JSON(200, result)
		},
	}

	a.peopleComponent.SearchHousehold(c, data, output).Wait()
}