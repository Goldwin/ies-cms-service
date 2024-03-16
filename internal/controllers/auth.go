package controllers

import (
	"context"
	"errors"
	"strings"

	"github.com/Goldwin/ies-pik-cms/internal/bus"
	"github.com/Goldwin/ies-pik-cms/internal/bus/common"
	"github.com/Goldwin/ies-pik-cms/internal/middleware"
	output "github.com/Goldwin/ies-pik-cms/internal/out/auth"
	"github.com/Goldwin/ies-pik-cms/pkg/auth"
	"github.com/Goldwin/ies-pik-cms/pkg/auth/dto"
	"github.com/Goldwin/ies-pik-cms/pkg/common/out"
	people "github.com/Goldwin/ies-pik-cms/pkg/people/dto"
	"github.com/gin-gonic/gin"
	"github.com/vmihailenco/msgpack/v5"
)

type authController struct {
	authComponent       auth.AuthComponent
	authOutputComponent output.AuthOutputComponent
}

func InitializeAuthController(r *gin.Engine, authComponent auth.AuthComponent,
	eventBusComponent bus.EventBusComponent,
	authOutputComponent output.AuthOutputComponent,
	middlewareComponent middleware.MiddlewareComponent,
) {
	authController := authController{
		authComponent:       authComponent,
		authOutputComponent: authOutputComponent,
	}

	authGroup := r.Group("auth")
	authGroup.GET("", authController.auth)
	authGroup.POST("registration", middlewareComponent.Auth(), authController.completeRegistration)
	authGroup.POST("otp", authController.otp)
	authGroup.POST("otp/signin", authController.otpSignIn)
	authGroup.POST("password/signin", authController.passwordSignIn)

	eventBusComponent.Subscribe("people.added", func(ctx context.Context, event common.Event) {
		var person people.Person
		err := msgpack.Unmarshal(event.Body, &person)
		if err != nil {
			return
		}
		authComponent.CompleteRegistration(ctx, dto.CompleteRegistrationInput{
			FirstName:  person.FirstName,
			MiddleName: person.MiddleName,
			LastName:   person.LastName,
			Email:      person.EmailAddress,
		}, authOutputComponent.RegistrationOutput())
	})
}

func (a *authController) completeRegistration(c *gin.Context) {
	var input dto.CompleteRegistrationInput
	authRaw, ok := c.Get("auth_data")

	if !ok {
		c.AbortWithStatusJSON(401, gin.H{
			"error": "Unauthorized",
		})
		return
	}
	c.BindJSON(&input)
	authData, ok := authRaw.(middleware.AuthData)

	if !ok {
		c.AbortWithStatusJSON(400, gin.H{
			"error": "Malformed Token",
		})
		return
	}

	input.Email = authData.Email

	output := &outputDecorator[dto.AuthData]{
		output: a.authOutputComponent.RegistrationOutput(),
		errFunction: func(err out.AppErrorDetail) {
			c.JSON(400, gin.H{
				"error": err,
			})
		},
		successFunc: func(result dto.AuthData) {
			c.JSON(200, gin.H{
				"data": result,
			})
		},
	}
	a.authComponent.CompleteRegistration(c, input, output)
}

func (a *authController) auth(c *gin.Context) {
	var input dto.AuthInput
	header := c.GetHeader("Authorization")
	token, err := extractBearerToken(header)

	if err != nil {
		c.AbortWithStatusJSON(401, gin.H{
			"error": gin.H{
				"code":    1,
				"message": "Unauthorized",
			},
		})
		return
	}

	input.Token = token

	output := &outputDecorator[dto.AuthData]{
		output: a.authOutputComponent.AuthOutput(),
		errFunction: func(err out.AppErrorDetail) {
			c.AbortWithStatusJSON(400, gin.H{
				"error": err,
			})
		},
		successFunc: func(result dto.AuthData) {
			c.Set("auth_data", result)
			c.JSON(200, gin.H{
				"data": result,
			})
		},
	}
	a.authComponent.Auth(c, input, output)
}

func (a *authController) otp(c *gin.Context) {
	var input dto.OtpInput
	c.BindJSON(&input)
	a.authComponent.GenerateOtp(c, input, a.authOutputComponent.OTPOutput())
	c.JSON(204, gin.H{})
}

func (a *authController) otpSignIn(c *gin.Context) {
	var input dto.SignInInput
	c.BindJSON(&input)
	input.Method = "otp"
	output := &outputDecorator[dto.SignInResult]{
		errFunction: func(err out.AppErrorDetail) {
			c.JSON(400, gin.H{
				"error": err,
			})
		},
		successFunc: func(result dto.SignInResult) {
			c.JSON(200, gin.H{
				"data": result,
			})
		},
	}
	a.authComponent.SignIn(c, input, output)
}

func (a *authController) passwordSignIn(c *gin.Context) {
	var input dto.SignInInput
	c.BindJSON(&input)
	input.Method = "password"
	output := &outputDecorator[dto.SignInResult]{
		output: a.authOutputComponent.SignInOutput(),
		errFunction: func(err out.AppErrorDetail) {
			c.JSON(400, gin.H{
				"error": err,
			})
		},
		successFunc: func(result dto.SignInResult) {
			c.JSON(200, gin.H{
				"data": result,
			})
		},
	}
	a.authComponent.SignIn(c, input, output)
}

func extractBearerToken(bearer string) (string, error) {
	if bearer == "" {
		return "", errors.New("bad header value given")
	}

	jwtToken := strings.Split(bearer, " ")
	if len(jwtToken) != 2 {
		return "", errors.New("incorrectly formatted authorization header")
	}

	return jwtToken[1], nil
}
