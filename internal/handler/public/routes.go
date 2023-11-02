package public

import (
	"net/http"

	"github.com/dwnGnL/testWork/internal/application"
	"github.com/dwnGnL/testWork/internal/config"
	"github.com/dwnGnL/testWork/lib/external"
	"github.com/dwnGnL/testWork/lib/goerrors"
	"github.com/gin-gonic/gin"
)

type Handler struct {
}

func newHandler(cfg *config.Config) *Handler {
	return &Handler{}
}

func GenRouting(r *gin.RouterGroup, cfg *config.Config) {
	handler := newHandler(cfg)
	r.POST("/process", handler.Process)
}

func (h Handler) Process(c *gin.Context) {
	app, err := application.GetAppFromRequest(c)
	if err != nil {
		goerrors.Log().Warn("fatal err: %w", err)
		c.AbortWithStatus(http.StatusBadGateway)
		return
	}

	app.GetExternal().ProcessBatch(external.Item{})
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
