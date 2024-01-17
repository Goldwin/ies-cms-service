package main

import (
	"net/http"

	"github.com/Goldwin/ies-pik-cms/internal/bus"
	"github.com/Goldwin/ies-pik-cms/internal/config"
	controller "github.com/Goldwin/ies-pik-cms/internal/controllers"
	peopleData "github.com/Goldwin/ies-pik-cms/internal/data/people"
	"github.com/Goldwin/ies-pik-cms/internal/infra"
	"github.com/Goldwin/ies-pik-cms/internal/middleware"

	"github.com/Goldwin/ies-pik-cms/pkg/people"
	"github.com/gin-gonic/gin"
)

func main() {
	config := config.LoadConfig("people")

	infraComponent := infra.NewInfraComponent(config.InfraConfig)
	dataLayerComponent := peopleData.NewPeopleDataLayerComponent(config.DataConfig["PEOPLE"], infraComponent)
	peopleManagementComponent := people.NewPeopleManagementComponent(dataLayerComponent)
	middlewareComponent := middleware.NewMiddlewareComponent(config.MiddlewareConfig)
	eventBusComponent := bus.Redis(infraComponent)

	r := gin.Default()
	r.Use(middlewareComponent.Cors())

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	controller.InitializePeopleManagementController(r, middlewareComponent, peopleManagementComponent, eventBusComponent)

	r.Run()
}
