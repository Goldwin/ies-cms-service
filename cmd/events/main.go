package main

import (
	"net/http"

	"github.com/Goldwin/ies-pik-cms/internal/bus"
	"github.com/Goldwin/ies-pik-cms/internal/config"
	"github.com/Goldwin/ies-pik-cms/internal/controllers"
	data "github.com/Goldwin/ies-pik-cms/internal/data/events"
	"github.com/Goldwin/ies-pik-cms/internal/infra"
	"github.com/Goldwin/ies-pik-cms/internal/middleware"
	"github.com/Goldwin/ies-pik-cms/pkg/events"

	"github.com/gin-gonic/gin"
)

func main() {
	config := config.LoadConfig("events")

	infraComponent := infra.NewInfraComponent(config.InfraConfig)
	dataLayerComponent := data.NewChurchEventDataLayerComponent(config.DataConfig["EVENTS"], infraComponent)
	churchEventComponent := events.NewChurchEventComponent(dataLayerComponent)
	middlewareComponent := middleware.NewMiddlewareComponent(config.MiddlewareConfig)
	eventBusComponent := bus.Redis(infraComponent)

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	controllers.InitializeEventsController(r, middlewareComponent, churchEventComponent, eventBusComponent)

	r.Run()
}
