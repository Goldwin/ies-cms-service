package controllers

type ControllerConfig struct {
	Port int    `yaml:"port" default:"3000"`
	Mode string `yaml:"mode"`
}
