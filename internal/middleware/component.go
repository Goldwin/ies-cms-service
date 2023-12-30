package middleware

import "github.com/gin-gonic/gin"

type MiddlewareComponent interface {
	Auth(scopes ...string) gin.HandlerFunc
}

type middlewareComponentImpl struct {
	config MiddlewareConfig
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
