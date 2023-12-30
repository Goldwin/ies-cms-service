package middleware

type MiddlewareConfig struct {
	AuthUrl string `yaml:"authUrl" default:"http://localhost:8081"`
}
