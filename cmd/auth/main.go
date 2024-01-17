package main

import (
	"fmt"
	"net/http"

	"github.com/Goldwin/ies-pik-cms/internal/bus"
	"github.com/Goldwin/ies-pik-cms/internal/config"
	controller "github.com/Goldwin/ies-pik-cms/internal/controllers"
	authData "github.com/Goldwin/ies-pik-cms/internal/data/auth"
	"github.com/Goldwin/ies-pik-cms/internal/infra"
	"github.com/Goldwin/ies-pik-cms/internal/middleware"
	out "github.com/Goldwin/ies-pik-cms/internal/out/auth"
	"github.com/Goldwin/ies-pik-cms/pkg/auth"
	"github.com/gin-gonic/gin"
)

func main() {
	config := config.LoadConfig("auth")

	infraComponent := infra.NewInfraComponent(config.InfraConfig)
	authDataLayer := authData.NewAuthDataLayerComponent(config.DataConfig["AUTH"], infraComponent)
	authComponent := auth.NewAuthComponent(authDataLayer, config.Secret)
	eventBus := bus.Redis(infraComponent)
	authOutputComponent := out.NewAuthOutputComponent(eventBus)
	middlewareComponent := middleware.NewMiddlewareComponent(config.MiddlewareConfig)

	gin.SetMode(config.ControllerConfig.Mode)
	r := gin.Default()
	r.Use(middlewareComponent.Cors())

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	authComponent.Start()
	controller.InitializeAuthController(r, authComponent, eventBus, authOutputComponent, middlewareComponent)

	r.Run(fmt.Sprintf(":%d", config.ControllerConfig.Port))
}
