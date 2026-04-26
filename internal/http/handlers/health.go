package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
	appName   string
	version   string
	startedAt time.Time
}

func NewHealthHandler(appName, version string, startedAt time.Time) *HealthHandler {
	return &HealthHandler{
		appName:   appName,
		version:   version,
		startedAt: startedAt,
	}
}

func (h *HealthHandler) Liveness(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status":     "ok",
		"service":    h.appName,
		"version":    h.version,
		"started_at": h.startedAt.Format(time.RFC3339),
	})
}

func (h *HealthHandler) Readiness(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "ready",
		"service": h.appName,
	})
}
