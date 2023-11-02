package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/dwnGnL/testWork/internal/application"
	"github.com/dwnGnL/testWork/internal/handler/public"

	"github.com/dwnGnL/testWork/lib/goerrors"

	"github.com/dwnGnL/testWork/internal/config"
	"github.com/gin-gonic/gin"
)

type GracefulStopFuncWithCtx func(ctx context.Context) error

func SetupHandlers(core application.Core, cfg *config.Config) GracefulStopFuncWithCtx {
	c := gin.New()

	c.Use(application.WithApp(core), application.WithCORS())
	apiv1 := c.Group("/api/v1")

	public := apiv1.Group("/public")
	generatePublicRouting(public, cfg)

	port := fmt.Sprint(cfg.ListenPort)
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: c,
	}
	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			goerrors.Log().Fatalf("listen: %s\n", err)
		}
	}()
	return srv.Shutdown
}

func generatePublicRouting(gE *gin.RouterGroup, cfg *config.Config) {
	public.GenRouting(gE, cfg)
}
