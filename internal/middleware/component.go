package middleware

import (
	"github.com/gin-gonic/gin"
)

type MiddlewareComponent interface {
	Auth(scopes ...string) gin.HandlerFunc
	Cors() gin.HandlerFunc
}

type middlewareComponentImpl struct {
	config MiddlewareConfig
}

// Cors implements MiddlewareComponent.
func (m *middlewareComponentImpl) Cors() gin.HandlerFunc {
	cors := corsMiddleware{
		enabled: m.config.Cors,
	}
	return cors.Cors
}

// Auth implements MiddlewareComponent.
func (m *middlewareComponentImpl) Auth(scopes ...string) gin.HandlerFunc {
	auth := authMiddleware{
		scopes:  scopes,
		authUrl: m.config.AuthUrl,
	}
	return auth.Auth
}

func NewMiddlewareComponent(config MiddlewareConfig) MiddlewareComponent {
	return &middlewareComponentImpl{
		config: config,
	}
}
