package application

import (
	"context"

	"github.com/dwnGnL/testWork/internal/service"
	"github.com/gin-gonic/gin"
)

type Core interface {
	GetExternal() service.ExternalService
}

func WithApp(app Core) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), ContextApp, app)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
