package controllers

type ControllerConfig struct {
	Port int    `env:"CONTROLLER_PORT" default:"3000"`
	Mode string `env:"CONTROLLER_MODE"`
}
