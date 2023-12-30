package main

import (
	"fmt"
	"net/http"

	"github.com/Goldwin/ies-pik-cms/internal/bus"
	"github.com/Goldwin/ies-pik-cms/internal/config"
	controller "github.com/Goldwin/ies-pik-cms/internal/controllers"
	authData "github.com/Goldwin/ies-pik-cms/internal/data/auth"
	"github.com/Goldwin/ies-pik-cms/internal/infra"
	out "github.com/Goldwin/ies-pik-cms/internal/out/auth"
	"github.com/Goldwin/ies-pik-cms/pkg/auth"
	"github.com/gin-gonic/gin"
)

func main() {
	config := config.LoadConfig("auth")

	infraComponent := infra.NewInfraComponent(config.InfraConfig)
	authDataLayer := authData.NewAuthDataLayerComponent(config.DataConfig, infraComponent)
	authComponent := auth.NewAuthComponent(authDataLayer, config.Secret)
	eventBus := bus.NewLocalEventBusComponent()
	authOutputComponent := out.NewAuthOutputComponent(eventBus)

	gin.SetMode(config.ControllerConfig.Mode)
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	authComponent.Start()
	controller.InitializeAuthController(r, authComponent, authOutputComponent)

	r.Run(fmt.Sprintf(":%d", config.ControllerConfig.Port))
}
