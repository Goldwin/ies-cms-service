package config

import (
	"fmt"
	"log"
	"os"

	controller "github.com/Goldwin/ies-pik-cms/internal/controllers"
	"github.com/Goldwin/ies-pik-cms/internal/data"
	"github.com/Goldwin/ies-pik-cms/internal/infra"
	"github.com/Goldwin/ies-pik-cms/internal/middleware"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Secret           []byte
	InfraConfig      infra.InfraConfig           `yaml:"infrastructure"`
	DataConfig       data.DataLayerConfig        `yaml:"datalayer"`
	ControllerConfig controller.ControllerConfig `yaml:"controller"`
	MiddlewareConfig middleware.MiddlewareConfig `yaml:"middleware"`
}

func LoadConfig(module string) Config {
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
