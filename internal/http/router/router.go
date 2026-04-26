package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/jluisv16/hcm-go/internal/config"
	employeehttp "github.com/jluisv16/hcm-go/internal/employees/interfaces/http"
	"github.com/jluisv16/hcm-go/internal/http/handlers"
)

func New(cfg config.Config, healthHandler *handlers.HealthHandler, employeeHandler *employeehttp.Handler) *gin.Engine {
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	router.GET("/healthz", healthHandler.Liveness)
	router.GET("/readyz", healthHandler.Readiness)
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	v1 := router.Group("/api/v1")
	v1.GET("/ping", handlers.Ping)
	employeeHandler.Register(v1)

	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"service": cfg.AppName,
			"status":  "running",
		})
	})

	return router
}
