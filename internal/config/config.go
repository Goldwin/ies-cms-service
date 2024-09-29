package config

import (
	"log"
	"os"
	"strings"

	controller "github.com/Goldwin/ies-pik-cms/internal/controllers"
	"github.com/Goldwin/ies-pik-cms/internal/data"
	"github.com/Goldwin/ies-pik-cms/internal/infra"
	"github.com/Goldwin/ies-pik-cms/internal/middleware"
	"github.com/caarlos0/env/v10"
)

type Config struct {
	Secret           []byte                          `env:"SECRET"`
	ServiceName      string                          `env:"SERVICE_NAME,expand"`
	InfraConfig      infra.InfraConfig               `yaml:"infrastructure"`
	DataConfig       map[string]data.DataLayerConfig `yaml:"datalayer"`
	ControllerConfig controller.ControllerConfig     `yaml:"controller"`
	MiddlewareConfig middleware.MiddlewareConfig     `yaml:"middleware"`
	EmailConfig      infra.EmailConfig
}

func LoadConfigEnv() Config {
	config := Config{
		Secret: []byte(os.Getenv("SECRET_KEY")),
	}
	serviceName := os.Getenv("SERVICE_NAME")

	if err := env.Parse(&config); err != nil {
		log.Fatal(err.Error())
	}

	if err := env.Parse(&config.EmailConfig); err != nil {
		log.Fatal(err.Error())
	}

	config.ServiceName = serviceName
	config.MiddlewareConfig.AuthUrl = os.Getenv("AUTH_SERVICE_URL")
	modules := os.Getenv("SERVICE_MODULES")
	moduleList := strings.Split(modules, ",")
	config.DataConfig = make(map[string]data.DataLayerConfig)
	for _, module := range moduleList {
		opts := env.Options{
			Prefix: module + "_QUERY_",
		}
		queryConfig := data.WorkerConfig{}
		if err := env.ParseWithOptions(&queryConfig, opts); err != nil {
			log.Fatal(err.Error())
		}
		opts = env.Options{
			Prefix: module + "_COMMAND_",
		}
		workerConfig := data.WorkerConfig{}
		if err := env.ParseWithOptions(&workerConfig, opts); err != nil {
			log.Fatal(err.Error())
		}
		config.DataConfig[module] = data.DataLayerConfig{
			CommandConfig: &workerConfig,
			QueryConfig:   &queryConfig,
		}
	}

	return config
}
