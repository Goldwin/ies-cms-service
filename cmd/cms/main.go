package main

import (
	"net/http"

	"github.com/Goldwin/ies-pik-cms/internal/bus"
	"github.com/Goldwin/ies-pik-cms/internal/config"
	controller "github.com/Goldwin/ies-pik-cms/internal/controllers"
	"github.com/Goldwin/ies-pik-cms/internal/data"
	authData "github.com/Goldwin/ies-pik-cms/internal/data/auth"
	eventData "github.com/Goldwin/ies-pik-cms/internal/data/events"
	peopleData "github.com/Goldwin/ies-pik-cms/internal/data/people"

	out "github.com/Goldwin/ies-pik-cms/internal/out/auth"

	"github.com/Goldwin/ies-pik-cms/internal/infra"
	"github.com/Goldwin/ies-pik-cms/internal/middleware"

	"github.com/Goldwin/ies-pik-cms/pkg/auth"
	"github.com/Goldwin/ies-pik-cms/pkg/events"
	"github.com/Goldwin/ies-pik-cms/pkg/people"
	"github.com/gin-gonic/gin"
)

// Monolithic API. This contains all modules as 1 service
func main() {
	config := config.LoadConfigEnv()

	emailClient := infra.NewEmailClient(config.EmailConfig)

	infraComponent := infra.NewInfraComponent(config.InfraConfig)
	peopleDataLayer := peopleData.NewPeopleDataLayerComponent(config.DataConfig["PEOPLE"], infraComponent)
	authDataLayer := authData.NewAuthDataLayerComponent(data.DataLayerConfig{
		CommandConfig: &data.WorkerConfig{
			Mode:           "redis",
			DB:             "",
			UseTransaction: true,
		},
		QueryConfig: &data.WorkerConfig{
			Mode:           "redis",
			DB:             "",
			UseTransaction: true,
		},
	}, infraComponent)
	eventDataLayer := eventData.NewChurchEventDataLayerComponent(config.DataConfig["EVENTS"], infraComponent)

	authComponent := auth.NewAuthComponent(authDataLayer, config.Secret)
	peopleManagementComponent := people.NewPeopleManagementComponent(peopleDataLayer)
	churchEventComponent := events.NewChurchEventComponent(eventDataLayer)

	middlewareComponent := middleware.NewMiddlewareComponent(config.MiddlewareConfig)
	eventBusComponent := bus.Local()

	authOutputComponent := out.NewAuthOutputComponent(emailClient, eventBusComponent)

	r := gin.Default()
	r.Use(middlewareComponent.Cors())

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	controller.InitializePeopleManagementController(r, middlewareComponent, peopleManagementComponent, eventBusComponent)
	controller.InitializeAuthController(r, authComponent, eventBusComponent, authOutputComponent, middlewareComponent)
	controller.InitializeEventsController(r, middlewareComponent, churchEventComponent, eventBusComponent)
	controller.InitializeCMSController(r, authComponent, peopleManagementComponent, middlewareComponent, emailClient)

	r.Run()
}
