package middleware

type MiddlewareConfig struct {
	AuthUrl string `env:"AUTH_SERVICE_URL" yaml:"authUrl" default:"http://localhost:8081"`
	Cors    bool   `env:"USE_CORS" yaml:"cors" default:"true"`
}
