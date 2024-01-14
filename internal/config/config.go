package config

import (
	"fmt"
	"log"
	"os"
	"strings"

	controller "github.com/Goldwin/ies-pik-cms/internal/controllers"
	"github.com/Goldwin/ies-pik-cms/internal/data"
	"github.com/Goldwin/ies-pik-cms/internal/infra"
	"github.com/Goldwin/ies-pik-cms/internal/middleware"
	"github.com/caarlos0/env/v10"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Secret           []byte
	ServiceName      string                          `env:"SERVICE_NAME,expand"`
	InfraConfig      infra.InfraConfig               `yaml:"infrastructure"`
	DataConfig       map[string]data.DataLayerConfig `yaml:"datalayer"`
	ControllerConfig controller.ControllerConfig     `yaml:"controller"`
	MiddlewareConfig middleware.MiddlewareConfig     `yaml:"middleware"`
}

func LoadConfigEnv() Config {
	config := Config{
		Secret: []byte("secret"),
	}
	envName := os.Getenv("env")
	serviceName := os.Getenv("SERVICE_NAME")
	if envName == "" {
		envName = "dev"
	}

	secretKeyFile := fmt.Sprintf("%s/configs/%s/secret.key", os.Getenv("ROOT_DIR"), envName)
	buf, err := os.ReadFile(secretKeyFile)

	if err != nil {
		log.Fatalf("Failed to parse key file: %v", err)
	}

	opts := env.Options{
		Prefix: serviceName + "_",
	}

	if err := env.ParseWithOptions(&config, opts); err != nil {
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

	config.Secret = buf
	return config
}

func LoadConfig(module string) Config {
	return LoadConfigYaml(module)
}

func LoadConfigYaml(module string) Config {
	config := Config{
		Secret: []byte("secret"),
	}
	env := os.Getenv("env")
	if env == "" {
		env = "dev"
	}

	gin.SetMode(gin.DebugMode)
	secretKeyFile := fmt.Sprintf("%s/configs/%s/secret.key", os.Getenv("ROOT_DIR"), env)
	yamlFile := fmt.Sprintf("%s/configs/%s/%s.yaml", os.Getenv("ROOT_DIR"), env, module)
	buf, err := os.ReadFile(secretKeyFile)
	if err != nil {
		log.Fatalf("Failed to parse key file: %v", err)
	}

	cfg, err := os.ReadFile(yamlFile)
	if err != nil {
		log.Fatalf("Failed to parse yaml file: %v", err)
	}
	yaml.Unmarshal(cfg, &config)

	config.Secret = buf

	return config
}
