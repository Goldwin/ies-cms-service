package middleware

import "github.com/gin-gonic/gin"

type MiddlewareComponent interface {
	Auth(scopes ...string) gin.HandlerFunc
}

type middlewareComponentImpl struct {
}

// Auth implements MiddlewareComponent.
func (*middlewareComponentImpl) Auth(scopes ...string) gin.HandlerFunc {
	auth := authMiddleware{
		scopes: scopes,
	}
	return auth.Auth
}

func NewMiddlewareComponent() MiddlewareComponent {
	return &middlewareComponentImpl{}
}
