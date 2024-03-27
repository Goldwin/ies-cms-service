package controllers

import (
	"github.com/Goldwin/ies-pik-cms/internal/infra"
	"github.com/Goldwin/ies-pik-cms/internal/middleware"
	"github.com/Goldwin/ies-pik-cms/internal/out/cms"
	"github.com/Goldwin/ies-pik-cms/pkg/auth"
	"github.com/Goldwin/ies-pik-cms/pkg/auth/dto"
	"github.com/Goldwin/ies-pik-cms/pkg/common/out"
	"github.com/Goldwin/ies-pik-cms/pkg/people"
	peopleDto "github.com/Goldwin/ies-pik-cms/pkg/people/dto"
	"github.com/Goldwin/ies-pik-cms/pkg/people/queries"
	"github.com/gin-gonic/gin"
)

type cmsController struct {
	peopleComponent     people.PeopleManagementComponent
	authComponent       auth.AuthComponent
	resetPasswordOutput out.Output[dto.PasswordResetCodeResult]
	otpOutput           out.Output[dto.OtpResult]
}

type LoginResult struct {
	AccessToken string        `json:"token"`
	AuthData    dto.AuthData  `json:"auth_info"`
	Profile     *BasicProfile `json:"profile"`
}

type BasicProfile struct {
	ProfilePictureUrl string `json:"profile_picture"`
	FirstName         string `json:"first_name"`
	MiddleName        string `json:"middle_name"`
	LastName          string `json:"last_name"`
}

type GeneratePasswordCodeInput struct {
	EmailAddress string `json:"email"`
}

// BFF for IES Apps
func InitializeCMSController(r *gin.Engine,
	authComponent auth.AuthComponent,
	peopleComponent people.PeopleManagementComponent,
	middlewareComponent middleware.MiddlewareComponent,
	emailClient infra.EmailClient,
) {
	c := &cmsController{
		peopleComponent:     peopleComponent,
		authComponent:       authComponent,
		resetPasswordOutput: cms.NewPasswordResetOutputHandler(emailClient),
		otpOutput:           cms.NewRegistrationOTPOutputHandler(emailClient),
	}
	appGroup := r.Group("app")
	appGroup.POST("login", c.login)
	appGroup.POST("password/forgot", c.forgotPassword)
	appGroup.POST("password/reset", c.resetPassword)

	appGroup.POST("register/otp", c.otp)
}

func (a *cmsController) otp(c *gin.Context) {
	var input dto.OtpInput
	c.BindJSON(&input)
	a.authComponent.GenerateOtp(c, input, a.otpOutput)
	c.JSON(204, gin.H{})
}

func (a *cmsController) resetPassword(ctx *gin.Context) {
	var input dto.PasswordResetInput
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
	a.authComponent.ResetPassword(ctx, input, outputDecorator)
}

func (a *cmsController) login(ctx *gin.Context) {
	var input dto.SignInInput
	result := &LoginResult{}
	ctx.Bind(&input)

	input.Method = "password"

	peopleOutput := &outputDecorator[queries.ViewPersonResult]{
		output: nil,
		errFunction: func(err out.AppErrorDetail) {
			if err.Code != 404 {
				ctx.JSON(err.Code, gin.H{
					"error": err.Message,
				})
				return
			}
			ctx.JSON(200, gin.H{
				"data": result,
			})

		},
		successFunc: func(person queries.ViewPersonResult) {
			result.Profile = &BasicProfile{
				ProfilePictureUrl: person.Data.ProfilePictureUrl,
				FirstName:         person.Data.FirstName,
				MiddleName:        person.Data.MiddleName,
				LastName:          person.Data.LastName,
			}
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
			result.AccessToken = res.AccessToken
			result.AuthData = res.AuthData
			a.peopleComponent.ViewPersonByEmail(ctx, queries.ViewPersonByEmailQuery{
				Email: input.Email,
			}, peopleOutput)
		},
	}
	a.authComponent.SignIn(ctx, input, loginOutput)
}

func (a *cmsController) forgotPassword(ctx *gin.Context) {
	var input GeneratePasswordCodeInput
	ctx.Bind(&input)
	a.authComponent.GenerateResetToken(ctx, input.EmailAddress, &outputDecorator[dto.PasswordResetCodeResult]{
		output: a.resetPasswordOutput,
		errFunction: func(err out.AppErrorDetail) {
			ctx.JSON(400, gin.H{
				"error": err,
			})
		},
		successFunc: func(token dto.PasswordResetCodeResult) {
			ctx.JSON(204, gin.H{})
		},
	})
}

func (a *cmsController) UpdateProfile(ctx *gin.Context) {
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
