package controllers

import (
	"github.com/Goldwin/ies-pik-cms/pkg/auth"
	"github.com/Goldwin/ies-pik-cms/pkg/auth/dto"
	"github.com/Goldwin/ies-pik-cms/pkg/common/out"
	"github.com/Goldwin/ies-pik-cms/pkg/people"
	peopleDto "github.com/Goldwin/ies-pik-cms/pkg/people/dto"
	"github.com/Goldwin/ies-pik-cms/pkg/people/queries"
	"github.com/gin-gonic/gin"
)

type appController struct {
	peopleComponent people.PeopleManagementComponent
	authComponent   auth.AuthComponent
}

type LoginResult struct {
	AccessToken string `json:"token"`
	Email       string `json:"email"`
	FirstName   string `json:"firstName"`
	MiddleName  string `json:"middleName"`
	LastName    string `json:"lastName"`
}

func InitializeAppController(r *gin.Engine) {
	r.GET("health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})
}

func (a *appController) SavePassword(ctx *gin.Context) {
	var input dto.PasswordInput
	ctx.Bind(&input)

	outputDecorator := &outputDecorator[dto.PasswordResult]{
		output: nil,
		errFunction: func(err out.AppErrorDetail) {
			ctx.JSON(400, gin.H{
				"error": err,
			})
		},
		successFunc: func(result dto.PasswordResult) {
			ctx.JSON(200, gin.H{
				"data": result,
			})
		},
	}
	a.authComponent.SavePassword(ctx, input, outputDecorator)
}

func (a *appController) Login(ctx *gin.Context) {
	var input dto.SignInInput
	result := &LoginResult{}
	ctx.Bind(&input)

	peopleOutput := &outputDecorator[queries.ViewPersonResult]{
		output: nil,
		errFunction: func(err out.AppErrorDetail) {
			ctx.JSON(400, gin.H{
				"error": err,
			})
		},
		successFunc: func(person queries.ViewPersonResult) {
			ctx.JSON(200, gin.H{
				"data": result,
			})
		},
	}

	loginOutput := &outputDecorator[dto.SignInResult]{
		output: nil,
		errFunction: func(err out.AppErrorDetail) {
			ctx.JSON(400, gin.H{
				"error": err,
			})
		},
		successFunc: func(res dto.SignInResult) {
			result.Email = res.AuthData.Email
			result.AccessToken = res.AccessToken
			a.peopleComponent.ViewPersonByEmail(ctx, queries.ViewPersonByEmailQuery{
				Email: input.Email,
			}, peopleOutput)
		},
	}
	a.authComponent.SignIn(ctx, input, loginOutput)
}

func (a *appController) UpdateProfile(ctx *gin.Context) {
	var input peopleDto.Person
	ctx.Bind(&input)
	a.peopleComponent.UpdatePerson(ctx, input, &outputDecorator[peopleDto.Person]{
		output: nil,
		errFunction: func(err out.AppErrorDetail) {
			ctx.JSON(400, gin.H{
				"error": err,
			})
		},
		successFunc: func(person peopleDto.Person) {
			ctx.JSON(200, gin.H{
				"data": person,
			})
		},
	})
}
